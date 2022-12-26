package main

// run with go run main.go
import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mandelbrot/andre/main/compute"
	"mandelbrot/andre/main/model"
	"net/http"
	"os"
)

func handleHello(w http.ResponseWriter, req *http.Request) {

	log.Println(req.Method, req.RequestURI)

	// Returns hello world! as a response
	fmt.Fprintln(w, "Hello world!")
}

func handleCompute(w http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.RequestURI)
	reqBody, _ := io.ReadAll(req.Body)
	var data model.Compute_data
	json.Unmarshal([]byte(reqBody), &data)
	fmt.Println(data)
	file_path := compute.Compute(data)
	output_image := image_to_base64(file_path)
	fmt.Fprint(w, output_image)
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

func image_to_base64(file_path string) string {
	bytes, err := os.ReadFile(file_path)
	if err != nil {
		log.Fatal(err)
	}
	var base64Encoding string
	base64Encoding += "data:image/png;base64,"
	base64Encoding += base64.StdEncoding.EncodeToString(bytes)
	//fmt.Println(base64Encoding)
	return base64Encoding
}
