package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		name := os.Getenv("GO_SERVICE_NAME")
		if name == "" {
			name = "go-service"
		}
		fmt.Fprintf(w, "Hello from %s\n", name)
	})

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	log.Printf("Starting go-service on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
