package main

import (
	"log"
	"net/http"
)

func main() {

	// 1. Create a new http.ServeMux to route requests.
	ServeMux := http.NewServeMux()
	// 1.1 Add a handlers for the filesystem paths
	ServeMux.Handle("/", http.FileServer(http.Dir(".")))
	ServeMux.Handle("/assets", http.FileServer(http.Dir("/assets")))
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
