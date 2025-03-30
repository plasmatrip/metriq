// The main function is the entry point of the server application.
// It sets up a goroutine to listen for termination signals,
// sets up a logger and a storage object, and then starts the HTTP server.
// The server is configured with the routing package and a storage object.
// The storage object is an interface that provides methods to store and retrieve metrics.
// The server also starts a goroutine to perform backups of the storage object at regular intervals.
package main

import (
	"context"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"text/template"

	"google.golang.org/grpc"

	"github.com/plasmatrip/metriq/internal/backup"
	"github.com/plasmatrip/metriq/internal/logger"
	"github.com/plasmatrip/metriq/internal/server/config"
	srv "github.com/plasmatrip/metriq/internal/server/grpc"
	"github.com/plasmatrip/metriq/internal/server/router"
	"github.com/plasmatrip/metriq/internal/storage"
	"github.com/plasmatrip/metriq/internal/storage/db"
	"github.com/plasmatrip/metriq/internal/storage/mem"
	pb "github.com/plasmatrip/metriq/proto"
)

var buildVersion string
var buildDate string
var buildCommit string

const buildInfo = `
	Server build info
	Build version: {{if .BuildVersion}}{{.BuildVersion}}{{else}}"N/A"{{end}}
	Build date: {{if .BuildDate}}{{.BuildDate}}{{else}}"N/A"{{end}}
	Build commit: {{if .BuildCommit}}{{.BuildCommit}}{{else}}"N/A"{{end}}
	
`

// The main function sets up the server application and starts it.
// It sets up a goroutine to listen for termination signals and
// sets up a logger and a storage object. The storage object is an
// interface that provides methods to store and retrieve metrics.
// The logger is also an interface that provides methods to log messages.
// The server is configured with the routing package and a storage object.
// The server also starts a goroutine to perform backups of the storage object
// at regular intervals. The backup function takes a context, a storage object,
// a logger and a config object as arguments. The config object is used to
// determine the path to the backup file. The backup function is started in a
// goroutine and runs until the context is canceled. If the context is canceled,
// the backup function stops and the server is shut down.
func main() {
	// Create a context to listen for termination signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	templ := template.Must(template.New("buildInfo").Parse(buildInfo))

	data := struct {
		BuildVersion string
		BuildDate    string
		BuildCommit  string
	}{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
		BuildCommit:  buildCommit,
	}

	err := templ.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer log.Close()

	var stor storage.Repository
	if cfg.DSN == "" {
		stor = mem.NewStorage()
	} else {
		stor, err = db.NewPostgresStorage(ctx, cfg.DSN, log)
		if err != nil {
			log.Sugar.Infow("database connection error: ", err)
			return
		}
		defer stor.Close()
	}

	backup, err := backup.NewBackup(*cfg, stor, log)
	if err != nil {
		log.Sugar.Panic("error initializing backup: ", err, " ", cfg.FileStoragePath)
	}
	if cfg.DSN == "" {
		backup.Start(ctx)
	}

	var grpcServer *grpc.Server
	if cfg.EnableGRPC {
		listen, err := net.Listen("tcp", ":"+cfg.GRPCPort)
		if err != nil {
			log.Sugar.Panic("failed to listen: ", err)
		}
		grpcServer = grpc.NewServer()
		pb.RegisterMetricsServer(grpcServer, srv.NewGRPCServer(stor, *cfg, log))
		log.Sugar.Infow("The gRPC server is running. ", "Server address: ", cfg.GRPCPort)
		go func() {
			if errGrpc := grpcServer.Serve(listen); errGrpc != nil {
				log.Sugar.Info("Could not listen on tcp:"+cfg.GRPCPort+": ", errGrpc)
			}
		}()
	}

	server := http.Server{
		Addr: cfg.Host,
		Handler: func(next http.Handler) http.Handler {
			log.Sugar.Infow("The metrics collection server is running. ", "Server address: ", cfg.Host)
			log.Sugar.Infow("Server config", "store interval", cfg.StoreInterval, "backup file", cfg.FileStoragePath, "DSN", cfg.DSN, "KEY", cfg.Key)
			return next
		}(router.NewRouter(stor, *cfg, log)),
	}

	go server.ListenAndServe()

	// Wait for the context to be canceled
	<-ctx.Done()

	err = backup.Save()
	if err != nil {
		log.Sugar.Infow("error saving to backup: ", err, " ", cfg.FileStoragePath)
	}

	server.Shutdown(context.Background())

	if grpcServer != nil {
		grpcServer.GracefulStop()
	}

	log.Sugar.Infow("The server has been shut down gracefully")
}
