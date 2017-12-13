package api

import (
	"net/http"
)

//tokenAndIdentityAPIHandler is a wrapper over the mux router serving /token and /identities API
type tokenAndIdentityAPIHandler struct {
	apiRouter http.Handler
}

func (h *tokenAndIdentityAPIHandler) getRouter() http.Handler {
	return h.apiRouter
}

func (h *tokenAndIdentityAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.getRouter().ServeHTTP(w, r)
}

func NewTokenAndIdentityAPIHandler(clusterManagerCfg string, clusterCfg string, clusterName string) (http.Handler, error) {
	router, err := NewRouter(clusterManagerCfg, clusterCfg, clusterName)
	if err != nil {
		return nil, err
	}
	return &tokenAndIdentityAPIHandler{apiRouter: router}, nil
}
