package middleware

import (
	"net/http"

	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"go.uber.org/zap"

	"github.com/SGNL-ai/TraTs-Demo-Svcs/order/pkg/config"
)

func CombineMiddleware(middleware ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(middleware) - 1; i >= 0; i-- {
			final = middleware[i](final)
		}

		return final
	}
}

func GetMiddleware(orderConfig *config.OrderConfig, spireJwtSource *workloadapi.JWTSource, logger *zap.Logger) func(http.Handler) http.Handler {
	middlewareList := []func(http.Handler) http.Handler{}

	middlewareList = append(middlewareList, spiffeMiddleware(orderConfig, spireJwtSource, logger))

	return CombineMiddleware(middlewareList...)
}
