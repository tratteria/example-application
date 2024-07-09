package middleware

import (
	"net/http"
	"net/url"

	"github.com/coreos/go-oidc"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"golang.org/x/oauth2"

	"go.uber.org/zap"
)

func CombineMiddleware(middleware ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(middleware) - 1; i >= 0; i-- {
			final = middleware[i](final)
		}

		return final
	}
}

func GetMiddleware(oauth2Config oauth2.Config, oidcProvider *oidc.Provider, targetServiceSpiffeID spiffeid.ID, spiffeJwtSource *workloadapi.JWTSource, x509Source *workloadapi.X509Source, tratteriaURL *url.URL, tratteriaSpiffeID spiffeid.ID, traTToggle bool, logger *zap.Logger) func(http.Handler) http.Handler {
	middlewareList := []func(http.Handler) http.Handler{
		getAuthenticationMiddleware(oauth2Config, oidcProvider, logger),
	}

	if traTToggle {
		middlewareList = append(middlewareList, GetTxnTokenMiddleware(tratteriaURL, x509Source, tratteriaSpiffeID, logger))
	}

	middlewareList = append(middlewareList, getJwtSvidMiddleware(targetServiceSpiffeID, spiffeJwtSource, logger))

	return CombineMiddleware(middlewareList...)
}
