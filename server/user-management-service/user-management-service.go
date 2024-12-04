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

var port int = 8004

var db *sql.DB

// Data Structure
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
	router.HandleFunc("/api/v1/getUser/{id}", getUser).Methods("GET")
	router.HandleFunc("/api/v1/update/{id}", updateUser).Methods("PUT")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://127.0.0.1:8000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT"}),
	)(router)

	// Print port
	fmt.Printf("Listening at port %d\n", port)
	url := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(url, corsHandler))
}

// Get User
func getUser(w http.ResponseWriter, r *http.Request) {
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

	// Read from Database
	var user User
	query := "SELECT * FROM User WHERE UserID = ?"
	err = db.QueryRow(query, userID).Scan(&user.UserID, &user.Name, &user.EmailAddr, &user.ContactNo, &user.MembershipTier, &user.PasswordHash, &user.IsActivated, &user.VerificationCodeHash)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve user: %v", err), http.StatusInternalServerError)
		return
	}

	// Return User Data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

	w.WriteHeader(http.StatusAccepted)
}

// Update User Information
func updateUser(w http.ResponseWriter, r *http.Request) {
	// Connect to Database
	db, err := connectToDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error connecting to the database")
		return
	}

	// Read Data from Body
	var updatedUser User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Reads parameters
	params := mux.Vars(r)
	userID := params["id"]

	_, err = db.Exec(`
		UPDATE User
		SET Name = ?, ContactNo = ?, EmailAddr = ?
		WHERE UserID = ?`, updatedUser.Name, updatedUser.ContactNo, updatedUser.EmailAddr, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Something went wrong with updating")
		return
	}

	message := fmt.Sprintf("%s has been updated", updatedUser.Name)
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
