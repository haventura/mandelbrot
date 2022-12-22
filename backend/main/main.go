package main

// run with go run main.go
import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mandelbrot/andre/main/compute"
	"mandelbrot/andre/main/model"
	"net/http"
)

func handleHello(w http.ResponseWriter, req *http.Request) {

	log.Println(req.Method, req.RequestURI)

	// Returns hello world! as a response
	fmt.Fprintln(w, "Hello world!")
}

func handleCompute(w http.ResponseWriter, req *http.Request) {
	log.Println("it works")
	log.Println(req.Method, req.RequestURI)
	reqBody, _ := io.ReadAll(req.Body)
	var data model.Compute_data 
	json.Unmarshal(reqBody, &data)

	fmt.Println(data)
	

	compute.Compute(data)
	fmt.Fprintln(w, "Computed!")
}
	
func main() {
	// registers handleHello to GET /hello
	http.HandleFunc("/hello", handleHello)
	http.HandleFunc("/compute", handleCompute)
	// starts the server on port 5000
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatalln(err)
	}
}