package service

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"github.com/rancher/auth-proxy/model"
)

var router *mux.Router

//NewRouter creates and configures a mux router
func NewRouter() *mux.Router {
	router = mux.NewRouter().StrictSlash(true)

	// API framework routes

	// Application routes
	router.Methods("POST").Path("/v1/token").Handler(http.HandlerFunc(CreateToken))
	router.Methods("GET").Path("/v1/token").Handler(http.HandlerFunc(GetToken))
	router.Methods("DELETE").Path("/v1/token").Handler(http.HandlerFunc(DeleteToken))

	//router.Methods("POST").Path("/v1/config").Handler(api.ApiHandler(schemas, http.HandlerFunc(UpdateConfig)))
	//router.Methods("GET").Path("/v1/config").Handler(api.ApiHandler(schemas, http.HandlerFunc(GetConfig)))

	router.Methods("GET").Path("/v1/identities").Handler(http.HandlerFunc(GetIdentities))
	//router.Methods("GET").Path("/v1/identities").Handler(api.ApiHandler(schemas, http.HandlerFunc(SearchIdentities)))

	return router
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
