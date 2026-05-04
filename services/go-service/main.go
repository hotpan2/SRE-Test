package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var httpRequests = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "go_service_http_requests_total",
		Help: "Total number of HTTP requests received",
	})

func init() {
	prometheus.MustRegister(httpRequests)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	httpRequests.Inc()
	fmt.Fprint(w, "OK")
}

func greetingHandler(w http.ResponseWriter, r *http.Request) {
	httpRequests.Inc()
	name := os.Getenv("GO_SERVICE_NAME")
	if name == "" {
		name = "go-service"
	}
	fmt.Fprintf(w, "Hello from %s\n", name)
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/", greetingHandler)
	mux.Handle("/metrics", promhttp.Handler())
	return mux
}

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	log.Printf("Starting go-service on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, newMux()))
}
