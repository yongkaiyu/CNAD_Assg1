package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
	UserID         int       `json:"user_id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	PhoneNo        string    `json:"phone_no"`
	Password   	   string    `json:"password"`
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

/* func (m MembershipTier) IsValid() bool { To validate membership tier
    switch m {
    case MembershipTierBasic, MembershipTierPremium, MembershipTierVIP:
        return true
    default:
        return false
    }
} */

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/electric_car_sharing_db")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	router := mux.NewRouter()
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World in REST API!")
}
