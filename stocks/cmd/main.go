package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"go.uber.org/zap"

	"github.com/SGNL-ai/TraTs-Demo-Svcs/stocks/handler"
	"github.com/SGNL-ai/TraTs-Demo-Svcs/stocks/pkg/config"
	"github.com/SGNL-ai/TraTs-Demo-Svcs/stocks/pkg/database"
	"github.com/SGNL-ai/TraTs-Demo-Svcs/stocks/pkg/middleware"
	"github.com/SGNL-ai/TraTs-Demo-Svcs/stocks/pkg/service"
)

type App struct {
	Router         *mux.Router
	DB             *sql.DB
	Config         *config.StocksConfig
	SpireJwtSource *workloadapi.JWTSource
	HTTPClient     *http.Client
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
		logger.Fatal("Stocks database initialization failed.", zap.Error(err))
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
		Config:         appConfig,
		SpireJwtSource: spireJwtSource,
		HTTPClient:     &http.Client{},
		Logger:         logger,
	}

	middleware := middleware.GetMiddleware(appConfig, app.SpireJwtSource, app.HTTPClient, app.Logger)

	app.Router.Use(middleware)

	stockService := service.NewService(app.DB, app.Logger)
	stockHandler := handler.NewHandlers(stockService, app.Logger)

	app.initializeRoutes(stockHandler)

	srv := &http.Server{
		Handler:      app.Router,
		Addr:         "0.0.0.0:8070",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("Starting server on 8070.")
	log.Fatal(srv.ListenAndServe())
}

func (a *App) initializeRoutes(handlers *handler.Handlers) {
	a.Router.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")
	a.Router.HandleFunc("/api/stocks/search", handlers.SearchStocksHandler).Methods("GET")
	a.Router.HandleFunc("/api/stocks/holdings", handlers.GetUserHoldingsHandler).Methods("GET")
	a.Router.HandleFunc("/api/stocks/details/{id}", handlers.GetStockDetailsHandler).Methods("GET")
	a.Router.HandleFunc("/internal/stocks", handlers.UpdateUserStockHandler).Methods("POST")
}
