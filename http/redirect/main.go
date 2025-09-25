package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

const (
	defaultUser      = "guest"
	defaultReturnURL = "singletopactivity://callback"
)

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user
	user := r.URL.Query().Get("user")
	if user == "" {
		user = defaultUser
	}

	// Get the return URL
	returnURL := r.URL.Query().Get("return_url")
	if returnURL == "" {
		returnURL = defaultReturnURL // DEFAULT
	}

	// Generate token
	generatedToken := fmt.Sprintf("%06d", rand.Intn(900000)+100000) // 100000..999999
	redirectSuccess := fmt.Sprintf("%s?success=true&user=%s&token=%s", returnURL, user, generatedToken)
	redirectError := fmt.Sprintf("%s?success=false&user=%s", returnURL, user)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf(`
            <html>
            <head>
                <meta name="viewport" content="width=device-width, initial-scale=1.0">
                <style>
                    body {
                        font-family: sans-serif;
                        padding: 16px;
                    }
                    button {
                        font-size: 18px;
                        padding: 12px;
                        margin: 8px 0;
                        width: 100%%;
                    }
                </style>
            </head>
            <body>
                <h1>Test Browser Redirect</h1>
                <p>User: %s</p>
                <p>Generated token: %s</p>
				<p>Success ReturnUrl: %s</p>
				<p>Failure ReturnUrl: %s</p>
                <button onclick="window.location='%s'">Success</button>
                <button onclick="window.location='%s'">Error</button>

                <h4>See BrowserSwitch</h4>
                <a href="https://github.com/braintree/browser-switch-android" target="_blank">
                    GitHub: braintree/browser-switch-android
                </a>
            </body>
            </html>
        `, user, generatedToken, redirectSuccess, redirectError, redirectSuccess, redirectError)))
}

// go run main.go
// Desktop: http://localhost:8080/redirect?user=herman
// Android Emulator: http://10.0.2.2:8080/redirect?user=herman
func main() {
	http.HandleFunc("/redirect", redirectHandler)
	log.Println("🚀 Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
