document.addEventListener('DOMContentLoaded', () => {
    // Get vehicle ID from query parameters
    const vehicleId = localStorage.getItem('vehicleId');

    // Get user ID from localStorage
    const userId = localStorage.getItem('userId');

    // Populate the form with user ID and vehicle ID
    document.getElementById('user-id').value = userId || 'Not logged in';
    document.getElementById('vehicle-id').value = vehicleId || 'Unknown vehicle';

    // Format date to ISO 8601
    const formatDateTime = (input) => {
        const date = new Date(input);
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');
        const hours = String(date.getHours()).padStart(2, '0');
        const minutes = String(date.getMinutes()).padStart(2, '0');
        const seconds = String(date.getSeconds()).padStart(2, '0');
        
        // Get the timezone offset in hours and minutes
        //const timezoneOffset = date.getTimezoneOffset();
        //const offsetSign = timezoneOffset > 0 ? '-' : '+';
        //const offsetHours = String(Math.abs(Math.floor(timezoneOffset / 60))).padStart(2, '0');
        //const offsetMinutes = String(Math.abs(timezoneOffset % 60)).padStart(2, '0');
        
        // Format in the required format
        return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}Z`;
    };
    

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

        // Convert userId and vehicleId to integers
        const parsedUserId = parseInt(userId, 10);
        const parsedVehicleId = parseInt(vehicleId, 10);

        if (isNaN(parsedUserId) || isNaN(parsedVehicleId)) {
            alert('Invalid ID values. Please check user and vehicle IDs.');
            return;
        }

        // Format times to "YYYY-MM-DD HH:MM:SS"
        const formattedStartTime = formatDateTime(startTime);
        const formattedEndTime = formatDateTime(endTime);

        // Create booking object
        const bookingData = {
            user_id: parsedUserId,
            vehicle_id: parsedVehicleId,
            start_time: formattedStartTime,
            end_time: formattedEndTime,
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
