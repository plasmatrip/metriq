// This package contains the request handlers for the server application.
// Handlers are responsible for processing all incoming requests to the server.
// They are initialized with a repository interface and a logger instance.
// The repository is used to store and retrieve metrics from the server and the logger
// is used to log all errors and other events encountered by the handlers.
package handlers

import (
	"github.com/plasmatrip/metriq/internal/logger"
	"github.com/plasmatrip/metriq/internal/server/config"
	"github.com/plasmatrip/metriq/internal/storage"
)

// The Handlers struct is a key component of the server application, responsible for handling
// all incoming HTTP requests. It is composed of three main fields: Repo, config, and lg.
//
// - Repo: This is an interface that represents the storage repository used by the handlers to
//   store and retrieve metrics. The repository abstracts the underlying storage mechanism,
//   allowing for flexibility in how data is persisted, whether in-memory, in a database, or
//   another storage solution.
//
// - config: This field holds the server configuration settings, encapsulated in a Config struct.
//   These settings dictate various operational parameters of the server, such as timeouts,
//   server address, and other configurable options that influence the behavior of the handlers.
//
// - lg: The logger instance used to log errors, warnings, and informational messages. Logging
//   is crucial for debugging and monitoring the application, providing insights into the flow
//   of requests and any issues that arise during execution.
//
// The Handlers struct is instantiated via the NewHandlers function, which requires a repository,
// configuration, and logger as parameters to initialize a new instance. The handlers utilize
// these components to effectively manage the lifecycle of HTTP requests, ensuring data integrity,
// adherence to configuration, and comprehensive logging throughout the application's runtime.

type Handlers struct {
	Repo   storage.Repository
	config config.Config
	lg     logger.Logger
}

func NewHandlers(repo storage.Repository, config config.Config, lg logger.Logger) *Handlers {
	return &Handlers{Repo: repo, config: config, lg: lg}
}
