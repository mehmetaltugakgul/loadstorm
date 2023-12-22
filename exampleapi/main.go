package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	requestCount int
	mutex        sync.Mutex
)

func main() {
	http.HandleFunc("/", handleRequest)

	port := 8080
	serverAddr := fmt.Sprintf(":%d", port)

	fmt.Printf("Server listening on port %d...\n", port)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	requestCount++
	totalRequests := requestCount
	mutex.Unlock()

	log.Printf("Received request %d - Method: %s, URL: %s\n", totalRequests, r.Method, r.URL.Path)
	fmt.Fprintf(w, "Hello, this is the local API. Request number: %d\n", totalRequests)
}
