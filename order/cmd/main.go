package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"go.uber.org/zap"

	"github.com/SGNL-ai/TraTs-Demo-Svcs/order/handler"
	"github.com/SGNL-ai/TraTs-Demo-Svcs/order/pkg/config"
	"github.com/SGNL-ai/TraTs-Demo-Svcs/order/pkg/database"
	"github.com/SGNL-ai/TraTs-Demo-Svcs/order/pkg/middleware"
	"github.com/SGNL-ai/TraTs-Demo-Svcs/order/pkg/service"
)

type App struct {
	Router         *mux.Router
	DB             *sql.DB
	HTTPClient     *http.Client
	Config         *config.OrderConfig
	SpireJwtSource *workloadapi.JWTSource
	Logger         *zap.Logger
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Cannot initialize Zap logger: %v.", err)
	}

	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("Error syncing logger: %v", err)
		}
	}()

	db, err := database.InitializeDB(logger)
	if err != nil {
		logger.Fatal("Order database initialization failed.", zap.Error(err))
	}

	defer db.Close()

	appConfig := config.GetAppConfig()

	spireJwtSource, err := config.GetSpireJwtSource()
	if err != nil {
		logger.Fatal("Unable to create SPIRE JWTSource for fetching JWT-SVIDs.", zap.Error(err))
	}

	logger.Info("Successfully created SPIRE JWTSource for fetching JWT-SVIDs.")

	defer spireJwtSource.Close()

	app := &App{
		Router:         mux.NewRouter(),
		DB:             db,
		HTTPClient:     &http.Client{},
		Config:         appConfig,
		SpireJwtSource: spireJwtSource,
		Logger:         logger,
	}

	middleware := middleware.GetMiddleware(app.Config, app.SpireJwtSource, app.Logger)

	app.Router.Use(middleware)

	orderService := service.NewService(app.DB, app.HTTPClient, app.Config, app.SpireJwtSource, app.Logger)
	orderHandler := handler.NewHandlers(orderService, app.Logger)

	app.initializeRoutes(orderHandler)

	srv := &http.Server{
		Handler:      app.Router,
		Addr:         "0.0.0.0:8090",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("Starting server on 8090.")
	log.Fatal(srv.ListenAndServe())
}

func (a *App) initializeRoutes(handlers *handler.Handlers) {
	a.Router.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")
	a.Router.HandleFunc("/api/order", handlers.OrderHandler).Methods("POST")
}
