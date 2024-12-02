package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var port int = 8001

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
	router.HandleFunc("/api/v1/cars", getAllCars).Methods("GET")
	router.HandleFunc("/api/v1/car/{id}", getCar).Methods("GET")
	router.HandleFunc("/api/v1/booking", postBooking).Methods("POST")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://127.0.0.1:8001"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT"}),
	)(router)

	// Print port
	fmt.Printf("Listening at port %d\n", port)
	url := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(url, corsHandler))
}

// Get All Cars
func getAllCars(w http.ResponseWriter, r *http.Request) {
	// Connect to Database
	db, err := connectToDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error connecting to the database")
		return
	}

	rows, err := db.Query("SELECT * FROM Car")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Errror executing DB query")
		return
	}
	defer rows.Close()

	var listOfCars []Car

	for rows.Next() {
		var c Car
		_ = rows.Scan(&c.CarID, &c.Model, &c.PlateNo, &c.RentalRate)
		listOfCars = append(listOfCars, c)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listOfCars)
	defer db.Close()
}

// Get Specific Car
func getCar(w http.ResponseWriter, r *http.Request) {
	// Connect to Database
	db, err := connectToDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error connecting to the database")
		return
	}

	// Reads parameters
	params := mux.Vars(r)
	carID := params["id"]

	// Read from Database
	var car Car
	query := "SELECT * FROM Car WHERE CarID = ?"
	err = db.QueryRow(query, carID).Scan(&car.CarID, &car.Model, &car.PlateNo, &car.RentalRate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve car: %v", err), http.StatusInternalServerError)
		return
	}

	// Return User Data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(car)

	w.WriteHeader(http.StatusAccepted)
}

// Create Booking
func postBooking(w http.ResponseWriter, r *http.Request) {
	// Connect to Database
	db, err := connectToDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error connecting to the database")
		return
	}

	// Read Data from Body
	var newBooking Booking
	err = json.NewDecoder(r.Body).Decode(&newBooking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Posting Booking into Database
	_, err = db.Exec(`
		INSERT INTO Booking (StartDate, EndDate, StartTime, EndTime, UserID, CarID, PaymentID)
		VALUES 
		(?, ?, ?, ?, ?, ?, ?)`, newBooking.StartDate, newBooking.EndDate, newBooking.StartTime, newBooking.EndTime, newBooking.UserID, newBooking.CarID, newBooking.PaymentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Something went wrong with creation")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Booking Posted Successfully!")
	w.WriteHeader(http.StatusOK)
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
