package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

const JWTSourceTimeout = 15 * time.Second

type spiffeIDs struct {
	Gateway spiffeid.ID
	Order   spiffeid.ID
	Stocks  spiffeid.ID
}

type StocksConfig struct {
	SpiffeIDs *spiffeIDs
}

func GetAppConfig() *StocksConfig {
	return &StocksConfig{
		SpiffeIDs: &spiffeIDs{
			Gateway: spiffeid.RequireFromString(getEnv("GATEWAY_SERVICE_SPIFFE_ID")),
			Order:   spiffeid.RequireFromString(getEnv("ORDER_SERVICE_SPIFFE_ID")),
			Stocks:  spiffeid.RequireFromString(getEnv("STOCKS_SERVICE_SPIFFE_ID")),
		},
	}
}

func GetSpireJwtSource() (*workloadapi.JWTSource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), JWTSourceTimeout)
	defer cancel()

	jwtSource, err := workloadapi.NewJWTSource(ctx)
	if err != nil {
		return nil, err
	}

	return jwtSource, nil
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		panic(fmt.Sprintf("%s environment variable not set", key))
	}

	return value
}
