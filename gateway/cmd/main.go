package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/SGNL-ai/TraTs-Demo-Svcs/gateway/handler"
	"github.com/SGNL-ai/TraTs-Demo-Svcs/gateway/pkg/config"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"go.uber.org/zap"
)

const SPIFEE_SOURCE_TIMEOUT = 15 * time.Second

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

	appConfig := config.GetAppConfig()

	oauth2Config := config.GetOauth2Config()
	oidcProvider := config.GetOIDCProvider(logger)

	jwtSourceCtx, cancel := context.WithTimeout(context.Background(), SPIFEE_SOURCE_TIMEOUT)
	defer cancel()

	spiffeJwtSource, err := workloadapi.NewJWTSource(jwtSourceCtx)
	if err != nil {
		logger.Fatal("Unable to create SPIFEE JWTSource for fetching JWT-SVIDs.", zap.Error(err))
	}

	defer spiffeJwtSource.Close()

	logger.Info("Successfully created SPIFEE JWTSource for fetching JWT-SVIDs.")

	x509SrcCtx, cancel := context.WithTimeout(context.Background(), SPIFEE_SOURCE_TIMEOUT)
	defer cancel()

	x509Source, err := workloadapi.NewX509Source(x509SrcCtx)
	if err != nil {
		logger.Fatal("Unable to create X509Source: " + err.Error())
	}

	defer x509Source.Close()

	logger.Info("Successfully created SPIFEE x509 source.")

	router := handler.SetupRoutes(appConfig, oauth2Config, oidcProvider, spiffeJwtSource, x509Source, logger)

	srv := &http.Server{
		Addr:         "0.0.0.0:30000",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("Starting server on 30000.")
	log.Fatal(srv.ListenAndServe())
}
