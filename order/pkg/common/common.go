package common

import (
	"net/url"
	"path"
)

func AppendPathToURL(baseURL *url.URL, appendPath string) *url.URL {
	newURL := *baseURL
	newURL.Path = path.Join(newURL.Path, appendPath)

	return &newURL
}

type ContextKey string

const TXN_TOKEN_CONTEXT_KEY ContextKey = "txn_token"