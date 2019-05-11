package main

import (
	"log"
	"net/http"

	"simio-api/resource"

	"github.com/gorilla/mux"
)

func main() {
	simioResource := resource.BuildSimioResource()
	router := mux.NewRouter()
	router.HandleFunc("/simian", simioResource.CheckSimian).Methods("POST")
	router.HandleFunc("/stats", simioResource.GetSimiansProportion).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
