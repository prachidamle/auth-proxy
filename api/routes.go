package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"github.com/rancher/auth-proxy/model"
)

var router *mux.Router

//NewRouter creates and configures a mux router
func NewRouter(clusterManagerCfg string, clusterCfg string, clusterName string) (*mux.Router, error) {

	apiServer, err := newAPIServer(clusterManagerCfg, clusterCfg, clusterName)
	if err != nil {
		return nil, err
	}

	router = mux.NewRouter().StrictSlash(true)
	// Application routes
	router.Methods("POST").Path("/v3/tokens").Queries("action", "login").Handler(http.HandlerFunc(apiServer.login))
	router.Methods("POST").Path("/v3/tokens").Queries("action", "logout").Handler(http.HandlerFunc(apiServer.logout))
	router.Methods("POST").Path("/v3/tokens").Handler(http.HandlerFunc(apiServer.deriveToken))
	router.Methods("GET").Path("/v3/tokens").Handler(http.HandlerFunc(apiServer.listTokens))

	router.Methods("GET").Path("/v3/identities").Handler(http.HandlerFunc(apiServer.listIdentities))
	//router.Methods("GET").Path("/v1/identities").Handler(api.ApiHandler(schemas, http.HandlerFunc(SearchIdentities)))

	return router, nil
}

//ReturnHTTPError handles sending out Error response
func ReturnHTTPError(w http.ResponseWriter, r *http.Request, httpStatus int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	err := model.AuthServiceError{
		Status:  strconv.Itoa(httpStatus),
		Message: errorMessage,
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(err)
}
