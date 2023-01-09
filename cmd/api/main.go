package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Danik14/library/internal/jsonlog"
	"github.com/Danik14/library/internal/models"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type config struct {
	port    int
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models models.Models
}

func main() {
	var cfg config

	// flag.IntVar(&cfg.port, "port", 4000, "API server port")
	// flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	// flag.Parse()
	// err := godotenv.Load("../../.env")
	// flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	// flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	// flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	// Create command line flags to read the setting values into the config struct.
	// Notice that we use true as the default for the 'enabled' setting?
	// flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	// flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	// flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	// flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelError)

	err := godotenv.Load()
	if err != nil {
		logger.PrintError(err, nil)
	}

	// creating a pointer to a mongodb struct
	var db *mongo.Client = DBSet(os.Getenv("DB-DSN"))

	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db.Database("library").Collection("users"),
			db.Database("library").Collection("books")),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")), //cfg.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.PrintInfo("Starting server on port: %s \n", map[string]string{
		"addr": srv.Addr,
		// "env": cfg.Env,
	})

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	// Making buffered channel, if not giving size
	// throwing warning, not sure if made perfectly
	quit := make(chan os.Signal, 1)

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")

}

func DBSet(dsn string) *mongo.Client {
	//connecting to local mongodb
	client, err := mongo.NewClient(options.Client().ApplyURI(dsn))
	if err != nil {
		log.Fatal(err)
	}

	//Set context with connection timeOut for secure connection
	dbContext, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err = client.Connect(dbContext)
	if err != nil {
		log.Println("Connection Failed")
		log.Fatal(err)
		return nil
	}

	fmt.Println("Successfuly Connected")
	return client
}
