package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var port int = 8001

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

	// Fetching db information
	godotenv.Load("./../.env")
	var dbUser = os.Getenv("DB_USER")
	var dbPassword = os.Getenv("DB_PASS")
	var dbName = os.Getenv("DB_NAME")

	// Constructing connection string
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", dbUser, dbPassword, dbName)

	// If not connected or there's an error, establish a new connection
	db, err := sql.Open("mysql", connectionString)
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
	router.HandleFunc("/api/v1/car/{id}", updateCarLocation).Methods("PUT")

	router.HandleFunc("/api/v1/booking", postBooking).Methods("POST")
	router.HandleFunc("/api/v1/checkValidity", checkBookingValidity).Methods("PUT")
	router.HandleFunc("/api/v1/booking/{id}", getBooking).Methods("GET")
	router.HandleFunc("/api/v1/booking/{id}", updateBooking).Methods("PUT")
	router.HandleFunc("/api/v1/booking/{id}", deleteBooking).Methods("DELETE")
	router.HandleFunc("/api/v1/booking/car/{id}", getBookingByCarID).Methods("GET")
	router.HandleFunc("/api/v1/booking/user/{id}", getBookingByUserID).Methods("GET")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://127.0.0.1:8001"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT"}),
	)(router)

	// Print port
	fmt.Printf("Listening at port %d\n", port)
	url := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(url, corsHandler))
}

//=========================== CAR RELATED ===========================

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
		_ = rows.Scan(&c.CarID, &c.Model, &c.PlateNo, &c.RentalRate, &c.Location)
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
	err = db.QueryRow(query, carID).Scan(&car.CarID, &car.Model, &car.PlateNo, &car.RentalRate, &car.Location)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve car: %v", err), http.StatusInternalServerError)
		return
	}

	// Return User Data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(car)

	w.WriteHeader(http.StatusAccepted)
}

// Update Car Location
func updateCarLocation(w http.ResponseWriter, r *http.Request) {
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

	// Read Data from Body
	var car Car
	err = json.NewDecoder(r.Body).Decode(&car)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Execute Query
	_, err = db.Exec(`
		UPDATE Car
		SET Location = ?
		WHERE CarID = ?`, car.Location, carID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Something went wrong with updating")
		return
	}

	message := fmt.Sprintf("Car %s's location has been successfully updated!", carID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)

	w.WriteHeader(http.StatusAccepted)
}

//=========================== BOOKING RELATED ===========================

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
		INSERT INTO Booking (Date, StartTime, EndTime, UserID, CarID, Model, PaymentID)
		VALUES 
		(?, ?, ?, ?, ?, ?, ?)`, newBooking.Date, newBooking.StartTime, newBooking.EndTime, newBooking.UserID, newBooking.CarID, newBooking.Model, newBooking.PaymentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Something went wrong with creation")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Booking Posted Successfully!")
	w.WriteHeader(http.StatusOK)
}

// Check Validity of Booking
func checkBookingValidity(w http.ResponseWriter, r *http.Request) {
	// Connect to Database
	db, err := connectToDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error connecting to the database")
		return
	}

	type BookingDetails struct {
		Date      string
		StartTime string
		EndTime   string
		CarID     int
	}

	type resultMessage struct {
		StatusCode    int
		ResultMessage string
	}

	// Read Data from Body
	var booking BookingDetails
	err = json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calling Stored Procedure
	var result resultMessage
	err = db.QueryRow(`
			CALL CheckBookingValidity(?, ?, ?, ?, @statusCode, @resultMessage)
			`, booking.Date, booking.StartTime, booking.EndTime, booking.CarID).Scan(&result.StatusCode, &result.ResultMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Something went wrong with checking")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
	w.WriteHeader(http.StatusOK)

}

// Get Specific Booking
func getBooking(w http.ResponseWriter, r *http.Request) {
	// Connect to Database
	db, err := connectToDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error connecting to the database")
		return
	}

	// Reads parameters
	params := mux.Vars(r)
	bookingID := params["id"]

	// Read from Database
	var b Booking
	query := "SELECT * FROM Booking WHERE BookingID = ?"
	err = db.QueryRow(query, bookingID).Scan(&b.BookingID, &b.Date, &b.StartTime, &b.EndTime, &b.UserID, &b.CarID, &b.Model, &b.PaymentID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve booking: %v", err), http.StatusInternalServerError)
		return
	}

	// Return User Data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(b)

	w.WriteHeader(http.StatusAccepted)
}

// List Bookings by CarID
func getBookingByCarID(w http.ResponseWriter, r *http.Request) {
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

	// Execute Query
	rows, err := db.Query("SELECT * FROM Booking WHERE CarID = ?", carID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Errror executing DB query")
		return
	}
	defer rows.Close()

	// Read Data
	var listOfBooking []Booking

	for rows.Next() {
		var b Booking
		_ = rows.Scan(&b.BookingID, &b.Date, &b.StartTime, &b.EndTime, &b.UserID, &b.CarID, &b.Model, &b.PaymentID)
		listOfBooking = append(listOfBooking, b)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listOfBooking)
	defer db.Close()
}

// List Bookings by UserID
func getBookingByUserID(w http.ResponseWriter, r *http.Request) {
	// Connect to Database
	db, err := connectToDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error connecting to the database")
		return
	}

	// Reads parameters
	params := mux.Vars(r)
	userID := params["id"]

	// Execute Query
	rows, err := db.Query("SELECT * FROM Booking WHERE UserID = ?", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Errror executing DB query")
		return
	}
	defer rows.Close()

	// Read Data
	var listOfBooking []Booking

	for rows.Next() {
		var b Booking
		_ = rows.Scan(&b.BookingID, &b.Date, &b.StartTime, &b.EndTime, &b.UserID, &b.CarID, &b.Model, &b.PaymentID)
		listOfBooking = append(listOfBooking, b)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listOfBooking)
	defer db.Close()
}

// Update Booking
func updateBooking(w http.ResponseWriter, r *http.Request) {
	// Connect to Database
	db, err := connectToDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error connecting to the database")
		return
	}

	// Reads parameters
	params := mux.Vars(r)
	bookingID := params["id"]

	// Read Data from Body
	var newBooking Booking
	err = json.NewDecoder(r.Body).Decode(&newBooking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Execute Query
	_, err = db.Exec(`
		UPDATE Booking
		SET Date = ?, StartTime = ?, EndTime = ?
		WHERE BookingID = ?`, newBooking.Date, newBooking.StartTime, newBooking.EndTime, bookingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Something went wrong with updating")
		return
	}

	message := fmt.Sprintf("Booking %s has been successfully updated!", bookingID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)

	w.WriteHeader(http.StatusAccepted)
}

// Delete Booking
func deleteBooking(w http.ResponseWriter, r *http.Request) {
	// Connect to Database
	db, err := connectToDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error connecting to the database")
		return
	}

	// Reads parameters
	params := mux.Vars(r)
	bookingID := params["id"]

	// Delete Booking
	_, err = db.Exec("DELETE FROM Booking WHERE BookingID = ?", bookingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Errror executing DB query")
		return
	}

	message := fmt.Sprintf("Booking %s has been deleted", bookingID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)

	w.WriteHeader(http.StatusAccepted)
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
