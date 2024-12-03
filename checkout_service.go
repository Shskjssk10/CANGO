package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"github.com/stripe/stripe-go/v81/tax/calculation"
	"github.com/stripe/stripe-go/v81/tax/transaction"
)

var port int = 8003

func main() {
	// Stripe Secret API KEY
	// Getting Secret Code
	godotenv.Load()
	stripe.Key = os.Getenv("STRIPE_KEY")

	router := mux.NewRouter()

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://127.0.0.1:8003"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT"}),
	)(router)

	// Routes
	router.HandleFunc("/api/v1/create-payment-intent", handleCreatePaymentIntent).Methods("POST")

	// Print port
	fmt.Printf("Listening at port %d\n", port)
	url := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(url, corsHandler))
}

type item struct {
	Id     string
	Amount int64
}

func calculateTax(items []item, currency stripe.Currency) *stripe.TaxCalculation {
	var lineItems []*stripe.TaxCalculationLineItemParams
	for _, item := range items {
		lineItems = append(lineItems, buildLineItem(item))
	}

	taxCalculationParams := &stripe.TaxCalculationParams{
		Currency: stripe.String(string(currency)),
		CustomerDetails: &stripe.TaxCalculationCustomerDetailsParams{
			Address: &stripe.AddressParams{
				Line1:      stripe.String("920 5th Ave"),
				City:       stripe.String("Seattle"),
				State:      stripe.String("WA"),
				PostalCode: stripe.String("98104"),
				Country:    stripe.String("US"),
			},
			AddressSource: stripe.String("shipping"),
		},
		LineItems: lineItems,
	}

	taxCalculation, _ := calculation.New(taxCalculationParams)
	return taxCalculation
}

func buildLineItem(i item) *stripe.TaxCalculationLineItemParams {
	return &stripe.TaxCalculationLineItemParams{
		Amount:    stripe.Int64(i.Amount), // Amount in cents
		Reference: stripe.String(i.Id),    // Unique reference for the item in the scope of the calculation
	}
}

// Securely calculate the order amount, including tax
func calculateOrderAmount(taxCalculation *stripe.TaxCalculation) int64 {
	// Calculate the order total with any exclusive taxes on the server to prevent
	// people from directly manipulating the amount on the client
	return taxCalculation.AmountTotal
}

func handleCreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Items []item `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewDecoder.Decode: %v", err)
		return
	}

	// Create a Tax Calculation for the items being sold
	taxCalculation := calculateTax(req.Items, "SGD")

	// Create a PaymentIntent with amount and currency
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(calculateOrderAmount(taxCalculation)),
		Currency: stripe.String(string(stripe.CurrencySGD)),
		// In the latest version of the API, specifying the `automatic_payment_methods` parameter is optional because Stripe enables its functionality by default.
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	params.AddMetadata("tax_calculation", taxCalculation.ID)

	pi, err := paymentintent.New(params)
	log.Printf("pi.New: %v", pi.ClientSecret)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("pi.New: %v", err)
		return
	}

	writeJSON(w, struct {
		ClientSecret   string `json:"clientSecret"`
		DpmCheckerLink string `json:"dpmCheckerLink"`
	}{
		ClientSecret: pi.ClientSecret,
		// [DEV]: For demo purposes only, you should avoid exposing the PaymentIntent ID in the client-side code.
		DpmCheckerLink: fmt.Sprintf("https://dashboard.stripe.com/settings/payment_methods/review?transaction_id=%s", pi.ID),
	})
}

// Invoke this method in your webhook handler when `payment_intent.succeeded` webhook is received
func handlePaymentIntentSucceeded(paymentIntent stripe.PaymentIntent) {
	// Create a Tax Transaction for the successful payment
	params := &stripe.TaxTransactionCreateFromCalculationParams{
		Calculation: stripe.String(paymentIntent.Metadata["tax_calculation"]),
		Reference:   stripe.String("myOrder_123"), // Replace with a unique reference from your checkout/order system
	}
	params.AddExpand("line_items")

	transaction.CreateFromCalculation(params)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewEncoder.Encode: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, &buf); err != nil {
		log.Printf("io.Copy: %v", err)
		return
	}
}
