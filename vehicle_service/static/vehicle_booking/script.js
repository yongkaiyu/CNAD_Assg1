document.addEventListener('DOMContentLoaded', () => {
    // Get vehicle ID from query parameters
    const vehicleId = localStorage.getItem('vehicleId');

    // Get user ID from localStorage
    const userId = localStorage.getItem('userId');

    // Populate the form with user ID and vehicle ID
    document.getElementById('user-id').value = userId || 'Not logged in';
    document.getElementById('vehicle-id').value = vehicleId || 'Unknown vehicle';

    // Handle form submission
    const bookingForm = document.getElementById('booking-form');
    bookingForm.addEventListener('submit', async (event) => {
        event.preventDefault();

        // Gather input values
        const startTime = document.getElementById('start-time').value;
        const endTime = document.getElementById('end-time').value;

        // Validate user and vehicle IDs
        if (!userId) {
            alert('User ID not found. Please log in.');
            return;
        }
        if (!vehicleId) {
            alert('Vehicle ID not found. Please select a vehicle.');
            return;
        }

        // Create booking object
        const bookingData = {
            user_id: userId,
            vehicle_id: vehicleId,
            start_time: startTime,
            end_time: endTime,
        };

        try {
            // Send booking data to the server
            const response = await fetch('/api/v1/booking/booking', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'userId': userId, // Include the userId from localStorage
                },
                body: JSON.stringify(bookingData),
            });

            console.log('Booking data being sent:', bookingData);

            const responseText = await response.text(); // Get raw response as text
            console.log('Raw server response:', responseText);

            // Parse JSON only if the response is successful
            if (response.ok) {
                const responseData = JSON.parse(responseText); // Convert to JSON
                alert('Booking confirmed!');
                window.location.href = '../bookings_home/';
            } else {
                const errorData = JSON.parse(responseText); // Handle error as JSON
                alert(`Booking failed: ${errorData.message}`);
            }
        } catch (error) {
            console.error('Error booking vehicle:', error);
            alert('An error occurred while booking the vehicle.');
        }
    });
});
