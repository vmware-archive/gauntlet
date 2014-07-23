package main

import (
	"fmt"
	"log"
	"net/http"
)

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	LastServerRequest = *r
}

var LastServerRequest http.Request

func StartWebServer() {

	http.HandleFunc("/results", resultsHandler)
	addr := "localhost:3000"
	fmt.Printf("listening on %s and responding to /results\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
