// The main function is the entry point of the server application.
// It sets up a goroutine to listen for termination signals,
// sets up a logger and a storage object, and then starts the HTTP server.
// The server is configured with the routing package and a storage object.
// The storage object is an interface that provides methods to store and retrieve metrics.
// The server also starts a goroutine to perform backups of the storage object at regular intervals.
package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"text/template"

	"github.com/plasmatrip/metriq/internal/backup"
	"github.com/plasmatrip/metriq/internal/logger"
	"github.com/plasmatrip/metriq/internal/server/config"
	"github.com/plasmatrip/metriq/internal/server/router"
	"github.com/plasmatrip/metriq/internal/storage"
	"github.com/plasmatrip/metriq/internal/storage/db"
	"github.com/plasmatrip/metriq/internal/storage/mem"
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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	t := template.Must(template.New("buildInfo").Parse(buildInfo))

	data := struct {
		BuildVersion string
		BuildDate    string
		BuildCommit  string
	}{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
		BuildCommit:  buildCommit,
	}

	err := t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}

	c, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	l, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer l.Close()

	var s storage.Repository
	if c.DSN == "" {
		s = mem.NewStorage()
	} else {
		s, err = db.NewPostgresStorage(ctx, c.DSN, l)
		if err != nil {
			l.Sugar.Infow("database connection error: ", err)
			return
			//os.Exit(1)
		}
		defer s.Close()
	}

	backup, err := backup.NewBackup(*c, s, l)
	if err != nil {
		l.Sugar.Panic("error initializing backup: ", err, " ", c.FileStoragePath)
	}
	if c.DSN == "" {
		backup.Start(ctx)
	}

	server := http.Server{
		Addr: c.Host,
		Handler: func(next http.Handler) http.Handler {
			l.Sugar.Infow("The metrics collection server is running. ", "Server address: ", c.Host)
			l.Sugar.Infow("Server config", "store interval", c.StoreInterval, "backup file", c.FileStoragePath, "DSN", c.DSN, "KEY", c.Key)
			return next
		}(router.NewRouter(s, *c, l)),
	}

	go server.ListenAndServe()

	<-ctx.Done()

	err = backup.Save()
	if err != nil {
		l.Sugar.Infow("error saving to backup: ", err, " ", c.FileStoragePath)
	}

	server.Shutdown(context.Background())

	l.Sugar.Infow("The server has been shut down gracefully")
}
