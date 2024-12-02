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
	"golang.org/x/crypto/bcrypt"
)

var port int = 8000

type User struct {
	UserID         int
	Name           string
	EmailAddr      string
	ContactNo      string
	MembershipTier string
	PasswordHash   string
}

var db *sql.DB

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

	// Routes
	router.HandleFunc("/api/v1/test", testingDB).Methods("GET")

	router.HandleFunc("/api/v1/getUser/{id}", getUser).Methods("GET")
	router.HandleFunc("/api/v1/registerUser", registerUser).Methods("POST")
	router.HandleFunc("/api/v1/loginUser", loginUser).Methods("GET")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://127.0.0.1:8000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT"}),
	)(router)

	fmt.Printf("Listening at port %d\n", port)
	url := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(url, corsHandler))
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	db, err := connectToDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error connecting to the database")
		return
	}

	var newUser User
	err = json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`
		INSERT INTO User (Name, EmailAddr, ContactNo, PasswordHash)
		VALUES 
		(?, ?, ?, ?)`, newUser.Name, newUser.EmailAddr, newUser.ContactNo, hashedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Something went wrong with creation")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("User Successully Created!")

	w.WriteHeader(http.StatusAccepted)
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string
		Password string
	}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := connectToDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error connecting to the database")
		return
	}

	var user User
	query := "SELECT * FROM User WHERE EmailAddr = ?"
	err = db.QueryRow(query, credentials.Email).Scan(&user.UserID, &user.Name, &user.EmailAddr, &user.ContactNo, &user.MembershipTier, &user.PasswordHash)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve user: %v", err), http.StatusInternalServerError)
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		fmt.Println("Login Unsuccessful")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
	fmt.Println("Login Unsuccessful")
}

func getUser(w http.ResponseWriter, r *http.Request) {
	// Connect to Database
	db, err := connectToDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error connecting to the database")
		return
	}

	params := mux.Vars(r)
	userID := params["id"]

	var user User
	query := "SELECT * FROM User WHERE UserID = ?"
	err = db.QueryRow(query, userID).Scan(&user.UserID, &user.Name, &user.EmailAddr, &user.ContactNo, &user.MembershipTier, &user.PasswordHash)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve user: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

	w.WriteHeader(http.StatusAccepted)
}

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
