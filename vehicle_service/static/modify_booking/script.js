// Populate form fields on page load
window.addEventListener('DOMContentLoaded', () => {
    const userId = localStorage.getItem('userId') || 'Not logged in';
    const bookingId = localStorage.getItem('bookingId') || 'Invalid booking';
    const vehicleId = localStorage.getItem('vehicleId') || 'Unknown vehicle';

    document.getElementById('user-id').value = userId;
    document.getElementById('booking-id').value = bookingId;
    document.getElementById('vehicle-id').value = vehicleId;
    
    // Handle form submission
    document.getElementById('modify-booking-form').addEventListener('submit', async function (e) {
        e.preventDefault();

        const startTime = document.getElementById('start-time').value;
        const endTime = document.getElementById('end-time').value;

        // Format date to ISO 8601
        const formatDateTime = (input) => {
            const date = new Date(input);
            const year = date.getFullYear();
            const month = String(date.getMonth() + 1).padStart(2, '0');
            const day = String(date.getDate()).padStart(2, '0');
            const hours = String(date.getHours()).padStart(2, '0');
            const minutes = String(date.getMinutes()).padStart(2, '0');
            const seconds = String(date.getSeconds()).padStart(2, '0');
            
            // Format in the required format
            return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
        };
    
        // Format times to "YYYY-MM-DD HH:MM:SS"
        const formattedStartTime = formatDateTime(startTime);
        const formattedEndTime = formatDateTime(endTime);
    
        if (!bookingId || !userId || !vehicleId) {
            alert("Booking ID, User ID, and Vehicle ID cannot be empty!");
            return;
        }
    
        try {
            const response = await fetch(`http://localhost:5000/api/v1/booking/modify/${bookingId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'userId': userId,
                    'vehicleId': vehicleId,
                },
                body: JSON.stringify({
                    startTime: formattedStartTime,
                    endTime: formattedEndTime,
                }),
            });
    
            if (response.ok) {
                alert('Booking updated successfully!');
                window.location.href = '../bookings_home/';
            } else {
                const error = await response.json();
                alert(`Error: ${error.message}`);
            }
        } catch (error) {
            console.error('Error:', error);
            alert('An error occurred while updating the booking.');
        }
    });
});