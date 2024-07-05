package common

import (
	"net/http"
	"net/url"
	"path"
)

type HttpMethod string

const (
	Get     HttpMethod = http.MethodGet
	Post    HttpMethod = http.MethodPost
	Put     HttpMethod = http.MethodPut
	Delete  HttpMethod = http.MethodDelete
	Patch   HttpMethod = http.MethodPatch
	Options HttpMethod = http.MethodOptions
)

var HttpMethodList = []HttpMethod{Get, Post, Put, Delete, Patch, Options}

type IDTokenClaims struct {
	Email string `json:"email"`
	Exp   int64  `json:"exp"`
}

type contextKey string

const OIDC_ID_TOKEN_CONTEXT_KEY contextKey = "id_token"

func AppendPathToURL(baseURL *url.URL, appendPath string) *url.URL {
	newURL := *baseURL
	newURL.Path = path.Join(newURL.Path, appendPath)
	return &newURL
}
