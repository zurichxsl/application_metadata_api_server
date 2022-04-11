package main

import (
	"application_metadata_api_server/server"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	log.Infof("Starting http httpServer...")
	httpServer := server.NewHttpServer()
	http.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})
	http.HandleFunc("/put", httpServer.PutHandler)
	http.HandleFunc("/get", httpServer.GetHandler)
	http.HandleFunc("/query", httpServer.SearchHandler)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
