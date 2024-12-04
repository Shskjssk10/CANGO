package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/MakMoinee/go-mith/pkg/email"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var port int = 8000

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
	router.HandleFunc("/api/v1/registerUser", registerUser).Methods("POST")
	router.HandleFunc("/api/v1/loginUser", loginUser).Methods("GET")
	router.HandleFunc("/api/v1/sendVerificationEmail", sendVerificationEmail).Methods("POST")
	router.HandleFunc("/api/v1/activateAccount", verifyVerificationCode).Methods("PUT")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://127.0.0.1:8000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT"}),
	)(router)

	// Print port
	fmt.Printf("Listening at port %d\n", port)
	url := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(url, corsHandler))
}

// Register User
func registerUser(w http.ResponseWriter, r *http.Request) {
	// Connect to Database
	db, err := connectToDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error connecting to the database")
		return
	}

	// Read Data from Body
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

	// Write into Database
	_, err = db.Exec(`
		INSERT INTO User (Name, EmailAddr, ContactNo, PasswordHash)
		VALUES 
		(?, ?, ?, ?)`, newUser.Name, newUser.EmailAddr, newUser.ContactNo, hashedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Something went wrong with creation")
		return
	}

	// Return Successful
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("User Successully Created!")

	w.WriteHeader(http.StatusAccepted)
}

// User Login
func loginUser(w http.ResponseWriter, r *http.Request) {
	// Read Data from Body
	var credentials struct {
		Email    string
		Password string
	}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connect to Database
	db, err := connectToDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error connecting to the database")
		return
	}

	// Retrieve User to be login
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
		// Return unsuccessful
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		fmt.Println("Login Unsuccessful")
		return
	}

	// Return successful
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
	fmt.Println("Login Unsuccessful")
}

// Send Email
func sendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	// Generating Code Hash
	var verificationCode string
	for i := 0; i < 6; i++ {
		num := rand.Intn(10)
		verificationCode += strconv.Itoa(num)
	}

	//Hash Verification Code
	hashVerificationCode, err := bcrypt.GenerateFromPassword([]byte(verificationCode), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Get User Email from Body
	var userEmail struct {
		Email string
	}

	err = json.NewDecoder(r.Body).Decode(&userEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// // Connect to Database
	db, err := connectToDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Update the user's record with the hashed verification code
	_, err = db.Exec("UPDATE User SET VerificationCodeHash = ? WHERE EmailAddr = ?", hashVerificationCode, userEmail.Email)
	if err != nil {
		http.Error(w, "Failed to update user record", http.StatusInternalServerError)
		return
	}

	// Getting Secret Code
	godotenv.Load("../.env")
	var emailKey = os.Getenv("EMAIL_KEY")

	// Send Email Verification Code
	emailService := email.NewEmailService(587, "smtp.gmail.com", "pookiebears2006@gmail.com", emailKey)

	isEmailSent, err := emailService.SendEmail(userEmail.Email, "Verification Email", fmt.Sprintf("Your verification code is: %s", verificationCode))
	if err != nil {
		log.Fatalf("Error sending email: %s", err)
	}

	if isEmailSent {
		log.Println("Email Send Successfully")
	} else {
		log.Println("Failed to send email")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Verification code sent successfully")
	w.WriteHeader(http.StatusOK)
}

// Verify Verification Code
func verifyVerificationCode(w http.ResponseWriter, r *http.Request) {
	// Read Data from Body
	var credentials struct {
		Email            string
		VerificationCode string
	}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connect to Database
	db, err := connectToDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Read from Database
	var user User
	query := "SELECT * FROM User WHERE EmailAddr = ?"
	err = db.QueryRow(query, credentials.Email).Scan(&user.UserID, &user.Name, &user.EmailAddr, &user.ContactNo, &user.MembershipTier, &user.PasswordHash, &user.IsActivated, &user.VerificationCodeHash)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve user: %v", err), http.StatusInternalServerError)
		return
	}

	// Compare verification code
	err = bcrypt.CompareHashAndPassword([]byte(user.VerificationCodeHash), []byte(credentials.VerificationCode))
	if err != nil {
		// Return unsuccessful
		http.Error(w, "Invalid verification code", http.StatusUnauthorized)
		fmt.Println("Verification Unsuccessful")
		return
	}

	// Update User Account is Verified
	_, err = db.Exec("UPDATE User SET IsActivated = 1 WHERE EmailAddr = ?", credentials.Email)
	if err != nil {
		http.Error(w, "Failed to update user record", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Account has been activiated!")
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
