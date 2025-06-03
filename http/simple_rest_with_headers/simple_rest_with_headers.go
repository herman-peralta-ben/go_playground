package main

// Mock Rest example

import (
	"fmt"
	"net/http"
)

// go run main.go
// test header: curl -H "X-API-KEY: secret123" http://localhost:8080/
func helloHandler(w http.ResponseWriter, r *http.Request) {
	const requiredKey = "secret123"
	apiKey := r.Header.Get("X-API-KEY")

	if apiKey != requiredKey {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Unauthorized: Missing or incorrect API key")
		return
	}

	fmt.Fprintln(w, "hello world")
}

func main() {
	http.HandleFunc("/", helloHandler)
	fmt.Println("Server listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
