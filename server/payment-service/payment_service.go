package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/MakMoinee/go-mith/pkg/email"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"github.com/stripe/stripe-go/v81/tax/calculation"
)

var port int = 8002

// Type Structures

type User struct {
	UserID               int
	Name                 string
	EmailAddr            string
	ContactNo            string
	MembershipTier       string
	PasswordHash         string
	IsActivated          int
	VerificationCodeHash string
}

type Car struct {
	CarID      int
	Model      string
	PlateNo    string
	RentalRate int
	Location   string
}

type Booking struct {
	BookingID int
	StartTime string
	EndTime   string
	Date      string
	CarID     int
	Model     string
	UserID    int
	PaymentID int
}

type Payment struct {
	PaymentID   int
	Amount      int
	DateCreated string
	Status      string
	UserID      int
	CarID       int
}

type item struct {
	Id     string
	Amount int64
}

var db *sql.DB

// Function to connect Database -- MUST BE USED AT ALL CRUD FUNCTIONS
func connectToDB() (*sql.DB, error) {
	if db != nil {
		// Check if the database connection is already established
		err := db.Ping()
		if err == nil {
			return db, nil
		}
	}

	// If not connected or there's an error, establish a new connection
	db, err := sql.Open("mysql", "root:Shskjssk10!@tcp(127.0.0.1:3306)/CNADAssg1DB")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
		return nil, err
	}

	// Ping the database to ensure the connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database:", err)
		return nil, err
	}

	return db, nil
}

func main() {
	// Stripe Secret API KEY
	// Getting Secret Code
	godotenv.Load("./../.env")
	stripe.Key = os.Getenv("STRIPE_KEY")

	router := mux.NewRouter()

	// Test Initial Database Connection
	router.HandleFunc("/api/v1/test", testingDB).Methods("GET")

	// Routes
	router.HandleFunc("/api/v1/payment", postPayment).Methods("POST")
	router.HandleFunc("/api/v1/paymentConfirmation", sendReceipt).Methods("POST")
	router.HandleFunc("/api/v1/create-payment-intent", handleCreatePaymentIntent).Methods("POST")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://127.0.0.1:8002"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT"}),
	)(router)

	// Print port
	fmt.Printf("Listening at port %d\n", port)
	url := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(url, corsHandler))
}

//=========================== PAYMENT RELATED TO DATABASE ===========================

// Create Payment
func postPayment(w http.ResponseWriter, r *http.Request) {
	// Connect to Database
	db, err := connectToDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error connecting to the database")
		return
	}

	// Read Data from Body
	var newPayment Payment
	err = json.NewDecoder(r.Body).Decode(&newPayment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Posting Payment into Database
	result, err := db.Exec(`
		INSERT INTO Payment (Amount, Status, UserID, CarID)
		VALUES 
		(?, 'Successful', ?, ?)`, newPayment.Amount, newPayment.UserID, newPayment.CarID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Something went wrong with creation")
		return
	}

	// Get the last inserted ID
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error getting last inserted ID")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Create a response struct
	type PaymentResponse struct {
		PaymentID int `json:"PaymentID"`
	}

	// Create a new PaymentResponse instance
	response := PaymentResponse{
		PaymentID: int(lastInsertId),
	}

	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}

func sendReceipt(w http.ResponseWriter, r *http.Request) {
	// Get Payment Details from Body
	type Email struct {
		Name      string
		EmailAddr string
		Model     string
		Date      string
		StartTime string
		EndTime   string
		Amount    int
	}

	var emailDetails Email

	err := json.NewDecoder(r.Body).Decode(&emailDetails)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Getting Secret Code
	godotenv.Load("../../.env")
	var emailKey = os.Getenv("EMAIL_KEY")

	// Send Email Verification Code
	emailService := email.NewEmailService(587, "smtp.gmail.com", "pookiebears2006@gmail.com", emailKey)

	var messageBody string

	messageBody = fmt.Sprintf(`
Dear %s,
	
This email confirms your booking for %s on %s %s to %s and that payment of $%d has been made.

Thank you for trusting us! We hope you have a wonderful time!
	`, emailDetails.Name, emailDetails.Model, emailDetails.Date, emailDetails.StartTime, emailDetails.EndTime, emailDetails.Amount)

	isEmailSent, err := emailService.SendEmail(emailDetails.EmailAddr, "Payment Confirmed", messageBody)
	if err != nil {
		log.Fatalf("Error sending email: %s", err)
	}

	if isEmailSent {
		log.Println("Email Send Successfully")
	} else {
		log.Println("Failed to send email")
	}
}

//=========================== STRIPE RELATED ===========================

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

// Actual function used for creating payment intent
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

// Test Database Connection
func testingDB(w http.ResponseWriter, r *http.Request) {
	db, err := connectToDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error connecting to the database")
		return
	}
	fmt.Println("Database has been successfully connected!")
	defer db.Close()
}
