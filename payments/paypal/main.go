package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	// Local
	"local/shared"
)

const (
	port    = "8081"
	baseURL = "https://api-m.sandbox.paypal.com"
)

var (
	clientID     string
	clientSecret string
)

// region Requests
type CreateOrderRequest struct {
	ReturnUrlScheme string `json:"return_url_scheme"`
	ReturnUrlHost   string `json:"return_url_host"`
}

type OrderRequest struct {
	Intent             string                   `json:"intent"`
	PurchaseUnits      []map[string]interface{} `json:"purchase_units"`
	ApplicationContext map[string]string        `json:"application_context"`
}

// endregion Requests

// region Responses
// AuthResponse representa el token de PayPal
type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// endregion Responses

func getEnvValue(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Couldn't load '%s', make sure it exists on your environment variables or on .env", key)
	}
	return val
}

func getAccessToken() (string, error) {
	req, _ := http.NewRequest("POST", baseURL+"/v1/oauth2/token", bytes.NewBufferString("grant_type=client_credentials"))
	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var auth AuthResponse
	if err := json.Unmarshal(body, &auth); err != nil {
		return "", err
	}

	return auth.AccessToken, nil
}

// Creates a Paypal Payment order (similar to Stripe's Payment Intent).
// Returns the orderID and approve links
func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	accessToken, err := getAccessToken()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var request CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("➡️ [/create-order] Request decode error: %v", err)
		http.Error(w, "invalid request", 400)
		return
	}

	log.Printf("➡️ [/create-order] Request return_url_scheme: %s, return_url_host: %s", request.ReturnUrlHost, request.ReturnUrlHost)

	order := OrderRequest{
		Intent: "CAPTURE",
		PurchaseUnits: []map[string]interface{}{
			{
				"amount": map[string]string{
					"currency_code": "USD",
					"value":         "10.00",
				},
			},
		},
		ApplicationContext: map[string]string{
			// 💡 notice that success is a query parameter
			"return_url": fmt.Sprintf("%s://%s?success=true", request.ReturnUrlScheme, request.ReturnUrlHost),
			"cancel_url": fmt.Sprintf("%s://%s?success=false", request.ReturnUrlScheme, request.ReturnUrlHost),
		},
	}

	body, _ := json.Marshal(order)

	req, _ := http.NewRequest("POST", baseURL+"/v2/checkout/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer res.Body.Close()

	respBody, _ := io.ReadAll(res.Body)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBody)
}

// Captures a previously payment already approved by the user.
func captureOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("orderId")
	if orderID == "" {
		http.Error(w, "missing orderId", 400)
		return
	}

	accessToken, err := getAccessToken()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v2/checkout/orders/%s/capture", baseURL, orderID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer res.Body.Close()

	respBody, _ := io.ReadAll(res.Body)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBody)
}

// ⚠️ Before running:
//
//	💡 Make sure you have a .env file containing PAYPAL_CLIENT_ID and PAYPAL_CLIENT_SECRET,
//	get them from https://developer.paypal.com/dashboard/applications/sandbox
//	💡 Use a personal test account: https://developer.paypal.com/dashboard/accounts
//
// $ go run .
func main() {
	log.Printf("🚀 Paypal Server running on http://localhost:%s", port)
	log.Printf("   🤖 Use http://10.0.2.2:%s/<api> on Android emulator", port)

	shared.LoadDotEnv(".env")

	clientID = getEnvValue("PAYPAL_CLIENT_ID")
	clientSecret = getEnvValue("PAYPAL_CLIENT_SECRET")

	http.HandleFunc("/create-order", createOrderHandler)
	http.HandleFunc("/capture-order", captureOrderHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
