package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
}

type Booking struct {
	BookingID int
	StartTime string
	EndTime   string
	StartDate string
	EndDate   string
	CarID     int
	UserID    int
	PaymentID int
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
	router := mux.NewRouter()

	// Test Initial Database Connection
	router.HandleFunc("/api/v1/test", testingDB).Methods("GET")

	// Routes

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://127.0.0.1:8002"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT"}),
	)(router)

	// Print port
	fmt.Printf("Listening at port %d\n", port)
	url := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(url, corsHandler))
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
