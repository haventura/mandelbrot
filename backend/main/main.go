package main

// run with go run main.go
// or  docker-compose --compatibility up --build
import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mandelbrot/andre/main/compute"
	"mandelbrot/andre/main/model"
	"net/http"
)

func handle_hello(w http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.RequestURI)
	fmt.Fprintln(w, "Hello world!")
}

func handle_compute_image(w http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.RequestURI)
	reqBody, _ := io.ReadAll(req.Body)
	var data model.Image_data
	json.Unmarshal([]byte(reqBody), &data)
	fmt.Println(data)
	image := compute.Compute_image(data)
	output_image := bytes_to_base64(image)
	fmt.Fprint(w, output_image)
}

func main() {
	http.HandleFunc("/hello", handle_hello)
	http.HandleFunc("/compute", handle_compute_image)
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatalln(err)
	}
}

func bytes_to_base64(data []byte) string {
	var base64Encoding string
	base64Encoding += "data:image/png;base64,"
	base64Encoding += base64.StdEncoding.EncodeToString(data)
	return base64Encoding
}
