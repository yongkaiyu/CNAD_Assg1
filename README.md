Header: CNAD Assg 1 Submission

go get -u github.com/gorilla/mux

go get -u "github.com/go-sql-driver/mysql"

go get -u "golang.org/x/crypto/bcrypt"

go get -u github.com/stripe/stripe-go

go get -u github.com/jung-kurt/gofpdf

Design consideration of microservices:

Architecture diagram:

Instructions for setting up and running microservices:

Download the file, add 3 users before setting up, to test membership_tier, require to manually alter in database

/* // Calculate the new duration and total amount
	startTime, err := time.Parse(time.RFC3339, input.StartTime)
	if err != nil {
		http.Error(w, "Invalid start time format", http.StatusBadRequest)
		log.Printf("Error parsing start time: %v", err)
		return
	}
	endTime, err := time.Parse(time.RFC3339, input.EndTime)
	if err != nil {
		http.Error(w, "Invalid end time format", http.StatusBadRequest)
		log.Printf("Error parsing end time: %v", err)
		return
	}

    duration := endTime.Sub(startTime).Hours()
	hours := int(math.Ceil(duration))
	if hours <= 0 {
		http.Error(w, "End time must be after start time", http.StatusBadRequest)
		return
	}

	fixedCost := 10.00
	totalAmount := float64(hours) * fixedCost

	result, err := db.Exec(`
        UPDATE bookings
        SET start_time = ?, end_time = ?
        WHERE booking_id = ? AND user_id = ?`,
		input.StartTime, input.EndTime, bookingID, userId)
	if err != nil {
		http.Error(w, "Error modifying booking", http.StatusInternalServerError)
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking update result", http.StatusInternalServerError)
		log.Printf("Error retrieving affected rows: %v", err)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "No changes made to the booking. Check input values.", http.StatusBadRequest)
		log.Printf("No rows updated for booking_id: %s, user_id: %s", bookingID, userId)
		return
	}

	// Update the billing record
	_, err = db.Exec(`
        UPDATE billings
        SET total_amount = ?
        WHERE booking_id = ?`,
		totalAmount, bookingID)
	if err != nil {
		http.Error(w, "Error updating billing entry", http.StatusInternalServerError)
		log.Printf("Error updating billing entry: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Booking modified successfully"}) */