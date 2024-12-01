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
	Phone          string    `json:"phone"`
	Password       string    `json:"password"`
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

	// Serve static files from "./static" at "/static/", allows localhost to serve html pages.
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// If any unmatched routes, redirected back here
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

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

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("JSON Decode Error: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Log the received data
	log.Printf("Received user data: %+v", user)

	name := user.Name
	email := user.Email
	phone := user.Phone
	password := user.Password
	membershipTier := "Basic"

	// Validate inputs
	if name == "" || email == "" || phone == "" || password == "" {
		log.Printf("Validation failed: name='%s', email='%s', phone='%s', password='%s'", name, email, phone, password)
		http.Error(w, "Name, email, phone and password are required", http.StatusBadRequest)
		return
	}

	// Encrypt password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insert user into database
	query := "INSERT INTO users (name, email, phone, password, membership_tier) VALUES (?, ?, ?, ?, ?)"
	_, err = db.Exec(query, name, email, phone, hashedPassword, membershipTier)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Prepare a JSON response
	response := map[string]string{
		"message": "User registered successfully",
	}

	// Set the response header to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Optional, defaults to 200
	json.NewEncoder(w).Encode(response)
}

func userAuthenticationHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var decodeUser User
	if err := json.NewDecoder(r.Body).Decode(&decodeUser); err != nil {
		log.Printf("JSON Decode Error: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	email := decodeUser.Email
	password := decodeUser.Password

	// Validate inputs
	if email == "" || password == "" {
		log.Printf("Validation failed: email='%s', password='%s'", email, password)
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Retrieve hashed password from database
	var hashedPassword string
	query := "SELECT password FROM users WHERE email = ?"
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

	// Retrieve user data
	var user User
	query2 := "SELECT user_id, name, email, phone, password FROM users WHERE email = ?"
	err2 := db.QueryRow(query2, email).Scan(&user.UserID, &user.Name, &user.Email, &user.Phone, &hashedPassword)
	if err2 != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	// Respond with user data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
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

	// Fetch benefits from membershipbenefits table
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

func httpError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]string{"error": message}
	json.NewEncoder(w).Encode(response)
}

func userProfileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // View Membership Status
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}
		var membershipTier string
		query := "SELECT membership_tier FROM users WHERE user_id = ?"
		err := db.QueryRow(query, userID).Scan(&membershipTier)
		if err != nil {
			if err == sql.ErrNoRows {
				httpError(w, "User not found", http.StatusNotFound)
			} else {
				httpError(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		// Respond with JSON data
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"membership_tier": membershipTier}
		json.NewEncoder(w).Encode(response)

		fmt.Fprintf(w, "Membership Tier: %s", membershipTier)

	case "PUT": // Update User Profile
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		var decodeUser User
		if err := json.NewDecoder(r.Body).Decode(&decodeUser); err != nil {
			log.Printf("JSON Decode Error: %v", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		name := decodeUser.Name
		email := decodeUser.Email
		phone := decodeUser.Phone
		password := decodeUser.Password

		// Validate inputs
		if name == "" || email == "" || phone == "" || password == "" {
			log.Printf("Validation failed: name='%s', email='%s', phone='%s', password='%s'", name, email, phone, password)
			http.Error(w, "Name, email, phone and password are required", http.StatusBadRequest)
			return
		}

		var hashedPassword string

		// Only hash the password if it's provided
		if password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Failed to hash password", http.StatusInternalServerError)
				return
			}
			hashedPassword = string(hash)
		}

		// Update user profile
		query := "UPDATE users SET name = ?, email = ?, phone = ?"
		if hashedPassword != "" {
			query += ", password = ?"
		}
		query += " WHERE user_id = ?"

		if hashedPassword != "" {
			// Include hashed password in the update if provided
			_, err := db.Exec(query, name, email, phone, hashedPassword, userID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to update profile: %v", err), http.StatusInternalServerError)
				return
			}
		} else {
			// If no password update, just update other fields
			_, err := db.Exec(query, name, email, phone, userID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to update profile: %v", err), http.StatusInternalServerError)
				return
			}
		}

		// Respond with success message
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"message": "Profile updated successfully"}
		json.NewEncoder(w).Encode(response)
	}
}

func rentalHistoryHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	w.Header().Set("Content-Type", "application/json")

	// Ensure userID is provided
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	query := `SELECT booking_id, user_id, vehicle_id, start_time, end_time, status, total_cost, created_at, updated_at
			  FROM bookings WHERE user_id = ? AND status = "Completed" ORDER BY updated_at DESC`
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

	var rentals []map[string]interface{}

	for rows.Next() {
		var bookingID, userID, vehicleID, status, totalCost string
		var startTime, endTime, createdAt, updatedAt time.Time

		if err := rows.Scan(&bookingID, &userID, &vehicleID, &startTime, &endTime, &status, &totalCost, &createdAt, &updatedAt); err != nil {
			http.Error(w, "Error scanning rental data", http.StatusInternalServerError)
			return
		}

		rental := map[string]interface{}{
			"booking_id": bookingID,
			"user_id":    userID,
			"vehicle_id": vehicleID,
			"start_time": startTime,
			"end_time":   endTime,
			"status":     status,
			"total_cost": totalCost,
			"created_at": createdAt,
			"updated_at": updatedAt,
		}

		rentals = append(rentals, rental)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating rental data", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rentals)
}
