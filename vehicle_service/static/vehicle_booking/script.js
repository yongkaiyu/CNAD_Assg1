document.addEventListener('DOMContentLoaded', () => {
    // Get vehicle ID from query parameters
    const urlParams = new URLSearchParams(window.location.search);
    const vehicleId = urlParams.get('vehicleId');

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
            const response = await fetch('/api/v1/booking', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(bookingData),
            });

            if (response.ok) {
                alert('Booking confirmed!');
                // Optionally, redirect to another page
                window.location.href = '../bookings_home';
            } else {
                const errorData = await response.json();
                alert(`Booking failed: ${errorData.message}`);
            }
        } catch (error) {
            console.error('Error booking vehicle:', error);
            alert('An error occurred while booking the vehicle.');
        }
    });
});
