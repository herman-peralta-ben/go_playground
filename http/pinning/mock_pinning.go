package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("✅ Got Request:", r.Method, r.URL.Path)
	fmt.Fprintln(w, `{"message": "Hello, HTTPS from Go!"}`)
}

// $ go run mock_pinning.go
// Test: curl -k https://localhost:8443
func main() {
	http.HandleFunc("/", handler)

	fmt.Println("Listening on https://localhost:8443")
	err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
	if err != nil {
		panic(err)
	}
}
