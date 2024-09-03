package config

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
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
	SpiffeIDs          *spiffeIDs
	TratVerifyEndpoint *url.URL
	EnableTrats        bool
}

func GetAppConfig() *StocksConfig {
	return &StocksConfig{
		SpiffeIDs: &spiffeIDs{
			Gateway: spiffeid.RequireFromString(getEnv("GATEWAY_SERVICE_SPIFFE_ID")),
			Order:   spiffeid.RequireFromString(getEnv("ORDER_SERVICE_SPIFFE_ID")),
			Stocks:  spiffeid.RequireFromString(getEnv("STOCKS_SERVICE_SPIFFE_ID")),
		},
		TratVerifyEndpoint: parseURL(getEnv("TRAT_VERIFY_ENDPOINT")),
		EnableTrats:        getBoolFromEnv("ENABLE_TRATS"),
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

func parseURL(rawurl string) *url.URL {
	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		panic(fmt.Sprintf("Error parsing URL %s: %v", rawurl, err))
	}

	return parsedURL
}

func getBoolFromEnv(key string) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic(fmt.Sprintf("%s environment variable not set", key))
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		panic(fmt.Sprintf("Invalid boolean value for %s: %s", key, value))
	}

	return parsed
}
