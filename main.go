package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-openapi/spec"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"gopkg.in/yaml.v2"
)

var (
	crdFilePath string
)

func init() {
	flag.StringVar(&crdFilePath, "crd-file-path", "", "Path to the Kubernetes CRD YAML file")
	flag.Parse()

	if crdFilePath == "" {
		log.Fatal("Please provide the path to the Kubernetes CRD YAML file using the -crd-file-path flag")
	}
}

func main() {
	// Load the Kubernetes CRD YAML file
	crdFile, err := ioutil.ReadFile(crdFilePath)
	if err != nil {
		log.Fatalf("Failed to read Kubernetes CRD YAML file: %v", err)
	}

	// Parse the YAML into a Swagger spec
	swagger := spec.Swagger{}
	err = yaml.Unmarshal(crdFile, &swagger)
	if err != nil {
		log.Fatalf("Failed to parse Kubernetes CRD YAML into Swagger spec: %v", err)
	}

	// Create a new Gorilla Mux router
	r := mux.NewRouter()

	// Serve the Swagger spec at the root URL
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := yaml.Marshal(&swagger)
		if err != nil {
			log.Fatalf("Failed to marshal Swagger spec: %v", err)
		}

		w.Header().Set("Content-Type", "application/yaml")
		_, err = w.Write(data)
		if err != nil {
			log.Fatalf("Failed to write Swagger spec: %v", err)
		}
	})

	// Serve the Swagger UI at the "/docs" URL
	uiHandler := http.StripPrefix("/docs", http.FileServer(http.Dir("swagger-ui/")))
	r.PathPrefix("/docs").Handler(alice.New().Then(uiHandler))

	// Start the server
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
