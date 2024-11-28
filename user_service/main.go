package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID         int       `json:"user_id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	PhoneNo        string    `json:"phone_no"`
	PasswordHash   string    `json:"password_hash"`
	MembershipTier string    `json:"membership_tier"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type MembershipBenefits struct {
	Tier           string  `json:"tier"`
	HourlyRate     float64 `json:"hourly_rate"`
	PriorityAccess bool    `json:"priority_access"`
	BookingLimit   int     `json:"booking_limit"`
}

type Booking struct {
	BookingID int       `json:"booking_id"`
	UserID    int       `json:"user_id"`
	VehicleID int       `json:"vehicle_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
	TotalCost *float64  `json:"total_cost,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const (
	StatusActive    = "Active"
	StatusCompleted = "Completed"
	StatusCancelled = "Cancelled"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/electric_car_sharing_db")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/user/signup", userRegistrationHandler)
	router.HandleFunc("/api/v1/user/login", userAuthenticationHandler)
	router.HandleFunc("/api/v1/user/settings", userProfileHandler)
	router.HandleFunc("/api/v1/user/benefits", membershipBenefitsHandler)
	router.HandleFunc("/api/v1/user/history", rentalHistoryHandler)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

// Hash password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Verify password
func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func userRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	password := r.FormValue("password")
	membershipTier := "Basic" // Default tier

	// Encrypt the password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insert user into database
	query := "INSERT INTO users (name, email, phone, password_hash, membership_tier) VALUES (?, ?, ?, ?, ?)"
	_, err = db.Exec(query, name, email, phone, hashedPassword, membershipTier)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User registered successfully")
}

func userAuthenticationHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// Retrieve hashed password from database
	var hashedPassword string
	query := "SELECT password_hash FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&hashedPassword)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare and verify passwords
	if !checkPassword(hashedPassword, password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Login successful")
}

func membershipBenefitsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse user_id from query parameters
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id parameter is required", http.StatusBadRequest)
		return
	}

	// Fetch membership tier for the user
	var membershipTier string
	query := "SELECT membership_tier FROM users WHERE user_id = ?"
	err := db.QueryRow(query, userID).Scan(&membershipTier)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		log.Printf("Error retrieving membership tier for user_id %s: %v", userID, err)
		return
	}

	// Fetch benefits from membership_benefits table
	var benefits MembershipBenefits
	query2 := "SELECT * FROM membershipbenefits WHERE tier = ?"
	err2 := db.QueryRow(query2, membershipTier).Scan(
		&benefits.Tier, &benefits.HourlyRate, &benefits.PriorityAccess, &benefits.BookingLimit)
	if err2 != nil {
		if err2 == sql.ErrNoRows {
			http.Error(w, "Membership tier not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		log.Printf("Error retrieving benefits for tier '%s': %v", membershipTier, err2)
		return
	}

	// Send benefits as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(benefits)

}
func userProfileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // View Membership Status
		userID := r.URL.Query().Get("user_id")
		var membershipTier string
		query := "SELECT membership_tier FROM users WHERE user_id = ?"
		err := db.QueryRow(query, userID).Scan(&membershipTier)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		fmt.Fprintf(w, "Membership Tier: %s", membershipTier)
	case "PUT": // Update User Profile
		userID := r.URL.Query().Get("user_id")
		name := r.FormValue("name")
		email := r.FormValue("email")
		phone := r.FormValue("phone")

		// Update user profile
		query := "UPDATE users SET name = ?, email = ?, phone = ? WHERE user_id = ?"
		_, err := db.Exec(query, name, email, phone, userID)
		if err != nil {
			http.Error(w, "Failed to update profile", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Profile updated successfully")
	}
}

func rentalHistoryHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	query := `SELECT * FROM bookings WHERE user_id = ? AND status = "Completed" ORDER BY updated_at DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		http.Error(w, "Error retrieving rental history", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Check if no rows found
	if !rows.Next() {
		http.Error(w, "No completed rentals found", http.StatusNotFound)
		return
	}

	// Prepare response
	fmt.Fprintf(w, "Rental History:\n")

	// Iterate through the rows and print booking details
	for rows.Next() {
		var bookingID, userID, vehicleID, status, totalCost string
		var startTime, endTime, createdAt, updatedAt time.Time

		// Scan booking details into variables
		if err := rows.Scan(&bookingID, &userID, &vehicleID, &startTime, &endTime, &status, &totalCost, &createdAt, &updatedAt); err != nil {
			http.Error(w, "Error scanning rental data", http.StatusInternalServerError)
			return
		}

		// Display the booking details
		fmt.Fprintf(w, "Booking ID: %s\n", bookingID)
		fmt.Fprintf(w, "User ID: %s\n", userID)
		fmt.Fprintf(w, "Vehicle ID: %s\n", vehicleID)
		fmt.Fprintf(w, "Start Time: %s\n", startTime)
		fmt.Fprintf(w, "End Time: %s\n", endTime)
		fmt.Fprintf(w, "Status: %s\n", status)
		fmt.Fprintf(w, "Total Cost: $%s\n\n", totalCost)
		fmt.Fprintf(w, "Created At: $%s\n\n", createdAt)
		fmt.Fprintf(w, "Updated At: $%s\n\n", updatedAt)
	}

	// Check for error from rows iteration
	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating rental data", http.StatusInternalServerError)
	}
}
