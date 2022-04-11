package server

import (
	"application_metadata_api_server/cache"
	"application_metadata_api_server/server/api"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// HttpServer is the interface of API server
type HttpServer interface {
	// PutHandler is the handler for put request, returns a new App Id
	PutHandler(w http.ResponseWriter, req *http.Request)
	// GetHandler is the handler for get request, returns the App content
	GetHandler(w http.ResponseWriter, req *http.Request)
	// SearchHandler is the handler for search request, returns a list of matching App Ids
	SearchHandler(w http.ResponseWriter, req *http.Request)
}

// httpServerImpl is an implementation of HttpServer
type httpServerImpl struct {
	store     cache.Store
	validator Validator
}

func NewHttpServer() HttpServer {
	return &httpServerImpl{
		store:     cache.InitStore(),
		validator: newAppValidator(),
	}
}

func (h *httpServerImpl) PutHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handleInternalError(w, err, "error reading request body")
		return
	}
	app, err := h.validator.ValidatePut(body)
	if err != nil {
		log.Warnf("invalid input yaml: %+v", err)
		handleValidationError(w, err)
		return
	}
	appId, err := h.store.Add(&app, body)
	if err != nil {
		handleInternalError(w, err, fmt.Sprintf("failed to put %+v", app))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "App Created"
	resp["id"] = string(appId)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("Error happened in JSON marshal error: %+v", err)
		handleInternalError(w, err, "json marshal error")
		return
	}
	w.Write(jsonResp)
	log.Infof("Successfully added app %s to the store", string(appId))
}

func (h *httpServerImpl) GetHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handleInternalError(w, err, "error reading request body")
		return
	}
	id := string(body)
	rawApp, err := h.store.Get(api.Id(id))
	if err != nil {
		handleNotFoundError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(rawApp)
}

func (h *httpServerImpl) SearchHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handleInternalError(w, err, "error reading request body")
		return
	}
	app, err := h.validator.ValidateSearch(body)
	if err != nil {
		log.Warnf("invalid input yaml: %+v", err)
		handleValidationError(w, err)
		return
	}
	rs, err := h.store.SearchStruct(&app)
	if err != nil {
		handleInternalError(w, err, fmt.Sprintf("failed to put %+v", app))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string][]api.Id)
	resp["result_list"] = rs
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("Error happened in JSON marshal error: %+v", err)
		handleInternalError(w, err, "json marshal error")
		return
	}
	w.Write(jsonResp)
	log.Infof("Found matched result %+v", rs)
}
