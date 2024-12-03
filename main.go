package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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

type Vehicle struct {
	VehicleID    int       `json:"vehicle_id"`
	LicensePlate string    `json:"license_plate"`
	Location     string    `json:"location"`
	ChargeLevel  string    `json:"charge_level"`
	Status       string    `json:"status"`
	Cleanliness  string    `json:"cleanliness"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

const (
	StatusAvailable   = "Available"
	StatusBooked      = "Booked"
	StatusMaintenance = "Maintenance"
)

const (
	CleanlinessClean    = "Clean"
	CleanlinessModerate = "Moderate"
	CleanlinessDirty    = "Dirty"
)

type BookedVehicle struct {
	BookingID    int    `json:"bookingId"`
	VehicleID    int    `json:"vehicleId"`
	LicensePlate string `json:"licensePlate"`
	Location     string `json:"location"`
	ChargeLevel  int    `json:"chargeLevel"`
	Status       string `json:"status"`
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
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

type Promotion struct {
	PromotionID        int       `json:"promotion_id"`
	Name               string    `json:"name"`
	DiscountPercentage float64   `json:"discount_percentage"`
	ExpiryDate         time.Time `json:"expiry_date"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type Billing struct {
	BillingID     int       `json:"billing_id"`
	BookingID     int       `json:"booking_id"`
	PaymentStatus string    `json:"payment_status"`
	PaymentMethod string    `json:"payment_method"`
	TotalAmount   float64   `json:"total_amount"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

const (
	PaymentStatusPending  = "Pending"
	PaymentStatusPaid     = "Paid"
	PaymentStatusRefunded = "Refunded"
)

const (
	PaymentMethodCreditCard = "Credit Card"
	PaymentMethodDebitCard  = "Debit Card"
	PaymentMethodPayPal     = "PayPal"
	PaymentMethodOther      = "Other"
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

	router.HandleFunc("/api/v1/booking/vehicles", availableVehiclesHandler)
	router.HandleFunc("/api/v1/booking/bookings", getBookedVehiclesHandler)
	router.HandleFunc("/api/v1/booking/booking", vehicleBookingHandler)
	router.HandleFunc("/api/v1/booking/modify", modifyBookingHandler)
	router.HandleFunc("/api/v1/booking/cancel", cancelBookingHandler)
	router.HandleFunc("/api/v1/booking/status", updateVehicleStatusHandler)

	// Serve static files from /static/{page}/ and route them to the corresponding service folder
	router.HandleFunc("/static/{page}/", serveStaticPage)

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./user_service/static/"))))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./vehicle_service/static/"))))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./billing_service/static/"))))

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

// Function to serve static pages dynamically from different services
func serveStaticPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"] // The page name (e.g., login, vehicles_available)

	var filePath string

	// Dynamically map to the correct folder
	switch page {
	case "login", "signup", "home", "settings", "history":
		filePath = "./user_service/static/" + page + "/index.html"
	case "vehicles_available", "vehicle_booking", "bookings_home":
		filePath = "./vehicle_service/static/" + page + "/index.html"
	case "billing_history", "payment":
		filePath = "./billing_service/static/" + page + "/index.html"
	default:
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	log.Printf("Serving file: %s", filePath)

	// Check if the file exists before serving
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("File not found: %s", filePath)
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	// Serve the file
	http.ServeFile(w, r, filePath)
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

		log.Printf("Connected to database: %v", db)
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

	// Ensure userID is provided
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	query := `SELECT booking_id, user_id, vehicle_id, start_time, end_time, status, total_cost, created_at, updated_at
			  FROM bookings WHERE user_id = ? AND status = "Completed" ORDER BY updated_at DESC`
	rows, err := db.Query(query, userIDInt)
	if err != nil {
		http.Error(w, "Error retrieving rental history", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var rentals []map[string]interface{}
	log.Println("Starting to iterate rows")
	log.Printf("Executing query with userID: %s", userID)

	for rows.Next() {
		var bookingID, userID, vehicleID, status string
		var startTimeStr, endTimeStr, createdAtStr, updatedAtStr string
		var totalCost float64

		if err := rows.Scan(&bookingID, &userID, &vehicleID, &startTimeStr, &endTimeStr, &status, &totalCost, &createdAtStr, &updatedAtStr); err != nil {
			log.Printf("Row scan error: %v", err)
			http.Error(w, "Error scanning rental data", http.StatusInternalServerError)
			return
		}

		// Parse start_time and end_time strings into time.Time
		startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr)
		if err != nil {
			log.Printf("Error parsing start_time: %v", err)
			http.Error(w, "Error parsing start time", http.StatusInternalServerError)
			return
		}

		endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
		if err != nil {
			log.Printf("Error parsing end_time: %v", err)
			http.Error(w, "Error parsing end time", http.StatusInternalServerError)
			return
		}

		// Parse created_at and updated_at strings into time.Time
		createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Printf("Error parsing created_at: %v", err)
			http.Error(w, "Error parsing created at time", http.StatusInternalServerError)
			return
		}

		updatedAt, err := time.Parse("2006-01-02 15:04:05", updatedAtStr)
		if err != nil {
			log.Printf("Error parsing updated_at: %v", err)
			http.Error(w, "Error parsing updated at time", http.StatusInternalServerError)
			return
		}

		// Format time.Time back into database string format since GMT +8 is added to time value
		formattedStartTime := startTime.Format("2006-01-02 15:04:05")
		formattedEndTime := endTime.Format("2006-01-02 15:04:05")
		formattedCreatedAt := createdAt.Format("2006-01-02 15:04:05")
		formattedUpdatedAt := updatedAt.Format("2006-01-02 15:04:05")

		rental := map[string]interface{}{
			"booking_id": bookingID,
			"user_id":    userID,
			"vehicle_id": vehicleID,
			"start_time": formattedStartTime,
			"end_time":   formattedEndTime,
			"status":     status,
			"total_cost": totalCost,
			"created_at": formattedCreatedAt,
			"updated_at": formattedUpdatedAt,
		}

		rentals = append(rentals, rental)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		http.Error(w, "Error iterating rental data", http.StatusInternalServerError)
		return
	}

	// Check if rentals is empty
	if len(rentals) == 0 {
		http.Error(w, "No completed rentals found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rentals)
}

func availableVehiclesHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query("SELECT vehicle_id, license_plate, location, charge_level, status, cleanliness, created_at, updated_at FROM vehicles WHERE status = 'Available' AND charge_level >= 20")
	if err != nil {
		http.Error(w, "Error fetching vehicles", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		var vehicle Vehicle
		err := rows.Scan(&vehicle.VehicleID, &vehicle.LicensePlate, &vehicle.Location, &vehicle.ChargeLevel, &vehicle.Status, &vehicle.Cleanliness, &vehicle.CreatedAt, &vehicle.UpdatedAt)
		if err != nil {
			http.Error(w, "Error scanning data", http.StatusInternalServerError)
			return
		}
		vehicles = append(vehicles, vehicle)
	}

	json.NewEncoder(w).Encode(vehicles)
}

func getBookedVehiclesHandler(w http.ResponseWriter, r *http.Request) {
	// Parse userId from query parameters (can be passed from localStorage on client-side)
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "userId parameter is required", http.StatusBadRequest)
		return
	}

	// Query the database to fetch vehicles booked by the user
	query := `
        SELECT 
            b.id AS bookingId,
            v.id AS vehicleId, 
            v.license_plate AS licensePlate, 
            v.location, 
            v.charge_level AS chargeLevel, 
            v.status, 
            b.start_time AS startTime, 
            b.end_time AS endTime 
        FROM 
            bookings b
        INNER JOIN 
            vehicles v ON b.vehicle_id = v.id
        WHERE 
            b.user_id = ?
    `
	rows, err := db.Query(query, userID)
	if err != nil {
		http.Error(w, "Error fetching booked vehicles", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Collect the results
	var bookedVehicles []BookedVehicle
	for rows.Next() {
		var vehicle BookedVehicle
		if err := rows.Scan(&vehicle.BookingID, &vehicle.VehicleID, &vehicle.LicensePlate, &vehicle.Location, &vehicle.ChargeLevel, &vehicle.Status, &vehicle.StartTime, &vehicle.EndTime); err != nil {
			http.Error(w, "Error scanning booked vehicles", http.StatusInternalServerError)
			return
		}
		bookedVehicles = append(bookedVehicles, vehicle)
	}

	// Convert results to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(bookedVehicles); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func vehicleBookingHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId") // Get userId from local storage (frontend sends this in headers)
	if userId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var booking Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(`
        INSERT INTO bookings (user_id, vehicle_id, start_time, end_time, total_cost)
        VALUES (?, ?, ?, ?, ?)`,
		userId, booking.VehicleID, booking.StartTime, booking.EndTime, booking.TotalCost)
	if err != nil {
		http.Error(w, "Error booking vehicle", http.StatusInternalServerError)
		return
	}

	// Update vehicle status to 'Booked'
	_, err = db.Exec(`UPDATE vehicles SET status = 'Booked' WHERE vehicle_id = ?`, booking.VehicleID)
	if err != nil {
		http.Error(w, "Error updating vehicle status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Vehicle booked successfully"})
}

func modifyBookingHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	var booking Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(`
        UPDATE bookings
        SET start_time = ?, end_time = ?, total_cost = ?
        WHERE booking_id = ? AND user_id = ?`,
		booking.StartTime, booking.EndTime, booking.TotalCost, booking.BookingID, userId)
	if err != nil {
		http.Error(w, "Error modifying booking", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Booking modified successfully"})
}

func cancelBookingHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	bookingID := r.URL.Query().Get("bookingId")

	_, err := db.Exec(`DELETE FROM bookings WHERE booking_id = ? AND user_id = ?`, bookingID, userId)
	if err != nil {
		http.Error(w, "Error canceling booking", http.StatusInternalServerError)
		return
	}

	// Optionally update vehicle status to 'Available'
	_, err = db.Exec(`UPDATE vehicles SET status = 'Available' WHERE id = (
        SELECT vehicle_id FROM bookings WHERE booking_id = ?
    )`, bookingID)
	if err != nil {
		http.Error(w, "Error updating vehicle status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Booking cancelled successfully"})
}

func updateVehicleStatusHandler(w http.ResponseWriter, r *http.Request) {
	var vehicle Vehicle
	if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(`
        UPDATE vehicles
        SET location = ?, charge_level = ?, cleanliness = ?
        WHERE vehicle_id = ?`,
		vehicle.Location, vehicle.ChargeLevel, vehicle.Cleanliness, vehicle.VehicleID)
	if err != nil {
		http.Error(w, "Error updating vehicle status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Vehicle status updated successfully"})
}
