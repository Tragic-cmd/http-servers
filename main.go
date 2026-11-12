package main

import (
	"log"
	"fmt"
	"net/http"
	"sync/atomic"
)

func main() {
	// Create a new http.ServeMux to route requests.
	ServeMux := http.NewServeMux()
	// Create a new apiConfig Struct
	apiCfg := apiConfig{}
	// Add a handlers for the filesystem paths
	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	ServeMux.Handle("/app/", apiCfg.middlewareMetricsInc(middlewareLog(fileServer)))
	ServeMux.Handle("/healthz", apiCfg.middlewareMetricsInc(middlewareLog(http.HandlerFunc(myHandler))))
	ServeMux.Handle("/metrics", http.HandlerFunc(apiCfg.apiHandler))
	ServeMux.Handle("/reset", http.HandlerFunc(apiCfg.resetMiddlewareMetrics))
	// Create a new http.Server struct.
	server := http.Server{
		Addr: ":8080",
		Handler: ServeMux,
	}
	// Start the server on port 8080, check for and log errors if they occur
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func myHandler(w http.ResponseWriter, r *http.Request) {
    // Set headers on w
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    // Set the status code on w
	w.WriteHeader(200)
    // Write data to w
	w.Write([]byte("OK"))
} 
func (cfg *apiConfig) apiHandler(w http.ResponseWriter, r *http.Request) {
    // Set headers on w
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    // Set the status code on w
	w.WriteHeader(200)
    // Write data to w
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
} 
func (cfg *apiConfig) resetMiddlewareMetrics(w http.ResponseWriter, r *http.Request) {
    // Set headers on w
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    // Set the status code on w
	w.WriteHeader(200)
	cfg.fileserverHits.Store(0)
    // Write data to w
	w.Write([]byte("Fileserver Hit Count Reset Successfully."))
} 

// Function that writes GET requests to the console log
func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
// Function that counts the number of requests made to the site
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}