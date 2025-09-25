package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
	"github.com/stripe/stripe-go/webhook"

	// Local
	"local/shared"
)

const port = "8080"

// run with `go run .`
func main() {
	// Set Stripe secret key
	shared.LoadDotEnv("../.env")
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Register routes
	http.HandleFunc("/create-one-click-checkout-card-payment-intent", handleCreateOneClickCheckoutCardPaymentIntent)
	http.HandleFunc("/create-unconfirmed-payment-intent", handleCreateUnconfirmedIntent)

	// TODO register webhook on Stripe
	http.HandleFunc("/webhook", handleStripeWebhook)

	log.Printf("🚀 Stripe Server running on http://localhost:%s", port)
	log.Printf("   🤖 Use http://10.0.2.2:%s/<api> on Android emulator", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type PaymentIntentRequest struct {
	PaymentMethodID string `json:"methodId"`
	Amount          int64  `json:"amount"`
	Currency        string `json:"currency"`
	UserId          string `json:"userId"`
	ProductId       string `json:"productId"`
}

type PaymentIntentResponse struct {
	ClientSecret string `json:"clientSecret"`
}

/*
* Confirmation
  - **One-click checkout** - Backend confirms (`Confirm: true`):
    Creates and confirms the PaymentIntent immediately  for simple payments without
	extra authentication. UI just finalizes the flow.
  - UI confirms (`Confirm: false`):
    Backend creates the PaymentIntent but UI confirms it using the `clientSecret`.
    Required for payments needing authentication (3D Secure, redirects like iDEAL, Klarna).

* Redirects
  - Redirecting payment methods include iDEAL (Dutch banks), Klarna (European banks), and 3D Secure.
  - Provide a `return_url` for payments with redirects (e.g. `AllowRedirects: stripe.String("always")`).
  - `return_url` is a deep link Stripe uses to redirect users back to your app after off-app payments.
  - In Flutter, set `returnURL` in `initPaymentSheet`.
  - In Android, configure an intent filter in `AndroidManifest.xml` for the URL scheme.
*/

func handleCreateOneClickCheckoutCardPaymentIntent(w http.ResponseWriter, r *http.Request) {
	var req PaymentIntentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("‼️ Invalid OneClickCheckoutCardPaymentRequest:", err)
		http.Error(w, "Invalid OneClickCheckoutCardPaymentRequest", http.StatusBadRequest)
		return
	}

	log.Printf("⬇️ OneClickCheckoutCardPaymentRequest:  {methodId=%s, amount=%d, currency=%s, userId=%s, productId=%s}\n",
		req.PaymentMethodID, req.Amount, req.Currency, req.UserId, req.ProductId)

	params := &stripe.PaymentIntentParams{
		PaymentMethod: stripe.String(req.PaymentMethodID),
		Amount:        stripe.Int64(req.Amount),    // Cents
		Currency:      stripe.String(req.Currency), // 3 letter code, e.g. "usd"
		//region confirm
		Confirm: stripe.Bool(true),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled:        stripe.Bool(true),
			AllowRedirects: stripe.String("never"), // do not redirect to external webpages
		},
		/*
			    AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			        Enabled:        stripe.Bool(true),
			        AllowRedirects: stripe.String("always"), // 👈 opcional, es el valor default
			    },
			    ReturnURL: stripe.String("myapp://stripe-redirect"),
				//
		*/
		//endregion confirm
		Metadata: map[string]string{
			"userId":    req.UserId,
			"productId": req.ProductId,
		},
	}

	intent, err := paymentintent.New(params)
	if err != nil {
		log.Println("‼️ Stripe error:", err)
		http.Error(w, "Failed to create payment intent", http.StatusInternalServerError)
		return
	}

	resp := PaymentIntentResponse{ClientSecret: intent.ClientSecret}

	log.Printf("⬆️ OneClickCheckoutCardPaymentResponse: %s\n", resp)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleCreateUnconfirmedIntent(w http.ResponseWriter, r *http.Request) {
	var req PaymentIntentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("‼️ Invalid UnconfirmedIntentRequest:", err)
		http.Error(w, "Invalid UnconfirmedIntentRequest", http.StatusBadRequest)
		return
	}

	log.Printf("⬇️ UnconfirmedIntentRequest:  {methodId=%s, amount=%d, currency=%s, userId=%s, productId=%s}\n",
		req.PaymentMethodID, req.Amount, req.Currency, req.UserId, req.ProductId)

	params := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(req.Amount),
		Currency:      stripe.String(req.Currency),
		PaymentMethod: stripe.String(req.PaymentMethodID),
		Confirm:       stripe.Bool(false), // Required to allow confirm later
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	intent, err := paymentintent.New(params)
	if err != nil {
		log.Printf("‼️ Stripe error: %v", err)
		http.Error(w, "Failed to create payment intent", http.StatusInternalServerError)
		return
	}

	resp := PaymentIntentResponse{
		ClientSecret: intent.ClientSecret,
	}

	log.Printf("⬆️ UnconfirmedIntentResponse: %s\n", resp)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	sigHeader := r.Header.Get("Stripe-Signature")
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		log.Printf("⚠️ Webhook signature verification failed: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if event.Type == "payment_intent.succeeded" {
		var intent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &intent); err != nil {
			log.Printf("Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userId := intent.Metadata["userId"]
		productId := intent.Metadata["productId"]

		log.Printf("✅💰 Payment succeeded for user %s, product %s", userId, productId)
		// TODO Save details to database
	} else {
		log.Printf("Unhandled event type: %s", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}
