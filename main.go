package main

import (
	"log"
	"net/http"
)

func main() {

	// 1. Create a new http.ServeMux to route requests.
	ServeMux := http.NewServeMux()
	// 1.1 Add a handlers for the filesystem paths
	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	ServeMux.Handle("/app/", middlewareLog(fileServer))
	ServeMux.Handle("/healthz", middlewareLog(http.HandlerFunc(myHandler)))
	// 2. Create a new http.Server struct.
	server := http.Server{
		Addr: ":8080",
		Handler: ServeMux,
	}
	// 3. Start the server on port 8080, check for and log errors if they occur
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

func myHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Set headers on w
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    // 2. Set the status code on w
	w.WriteHeader(200)
    // 3. Write data to w
	w.Write([]byte("OK"))
} 

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}