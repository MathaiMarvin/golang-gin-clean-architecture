// Package app configures and runs application.
package app

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/harmannkibue/golang_gin_clean_architecture/config"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/controller/http/v1"
	db "github.com/harmannkibue/golang_gin_clean_architecture/internal/entity/intfaces"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase/blog_usecase"
	"github.com/harmannkibue/golang_gin_clean_architecture/pkg/httpserver"
	"github.com/harmannkibue/golang_gin_clean_architecture/pkg/logger"
	"github.com/harmannkibue/golang_gin_clean_architecture/pkg/postgres"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
	"syscall"
)

// Run creates objects via constructors -.
func Run(cfg *config.Config) {

	//The line creates a new logger instance using the log level specified in the configuration. 
	//The logger is used to log messages to the console throughout the application.
	l := logger.New(cfg.Log.Level)

	// HTTP Server -.
	//the line below initializes a new instance of the gin Engine type and assigns it to the variable handler. 
	// THe new gin Engine instance already contains a default middleware attached. the middleware is for logging, handling static files and recovery from panics.
	handler := gin.Default()

	//This line establishes a connection to the POSTGRES database
	conn, err := postgres.New(cfg)

	if err != nil {
		fmt.Errorf("failed to connect to database %w", err)
	}

	// Defer closing the connection -.The defer keyword is used to schedule a function call ro be executed when the surrounding function returns. 
	//In this case, the defer keyword is used to schedule the closing of the database connection when the Run function returns. Inside the deffered function the close() method is called on the conn object to close db.
	//conn *sql.DB) This declares a parameter named conn of type *sql.DB.  The * before sql.DB indicates that conn is a pointer to a value of type sql.DB
	defer func(conn *sql.DB) {

		//In GO, the *sql.DB type represents a database connection pool. By using a pointer to sql.DB the code can pass the connection pool by reference to the deffered function, allowing it to close the connection pool when necessary.
		err := conn.Close()
		if err != nil {
			panic("ERROR CLOSING POSTGRES CONNECTION")
		}
	}(conn)

	// Initializing a store for repository -.The code below creates a new instance of a store object using the NewStore function from the db package. It assigns the newly created store object to the variable store. 
	store := db.NewStore(conn)

	//This line creates a new instance of the blog_usecase.BlogUseCase struct, which represents the use case for the blog functionality. 
	blogUsecase := blog_usecase.NewBlogUseCase(store, cfg)

	// Passing also the basic auth middleware to all  Routers - This line initializes the router for the version 1 API endpoints. It passes the gin router, logger, and the blog use case to the router.
	v1.NewRouter(handler, l, blogUsecase)

	// This line creates a new HTTP server instance using the gin router and port specified in the configuration. 
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal -. This line creates a buffered channel of type os.signal that is called interrupt. The buffer size is set to 1, allowing the channel to hold one signal at a time. 
	interrupt := make(chan os.Signal, 1)
	// This line registers the os.interrupt and syscall.SIGTERM signals that are sent to the interrupt channel. THis means that when the program receives an interrupt signal for instance (CTRL + C) it will be sent to the interrupt channel.
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()

	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
