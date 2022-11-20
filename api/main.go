package main

import (
	"fmt"
	"log"
	"net/http"
)

// handleHello GET /hello
func handleHello(w http.ResponseWriter, r *http.Request) {

	log.Println(r.Method, r.RequestURI)

	// Returns hello world! as a response
	fmt.Fprintln(w, "Hello world!")
}

// https://codeburst.io/load-balancing-go-api-with-docker-nginx-digital-ocean-d7f05f7c9b31
func main() {
	// registers handleHello to GET /hello
	http.HandleFunc("/hello", handleHello)
	// starts the server on port 5000
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatalln(err)
	}
}
