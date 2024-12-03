document.addEventListener('DOMContentLoaded', function() {
    // Get the userId from localStorage
    let userId = localStorage.getItem('userId');

    // Trim and sanitize userId so it wont be 1:1 but 1 instead
    if (userId) {
        userId = userId.trim().replace(/[^0-9]/g, ''); // Keep only numbers
    }

    if (!userId) {
        document.getElementById('rental-history').innerHTML = '<p class="error">User is not logged in.</p>';
        return;
    }

    // Fetch rental history from the API
    fetchRentalHistory(userId);
});

// Function to fetch rental history for the user
function fetchRentalHistory(userId) {
    fetch(`/api/v1/user/history?user_id=${userId}`)
        .then(response => response.json())
        .then(data => {
            if (data.length === 0) {
                document.getElementById('rental-history').innerHTML = '<p>No completed rentals found.</p>';
            } else {
                displayRentalHistory(data);
            }
        })
        .catch(error => {
            console.error('Error fetching rental history:', error);
            document.getElementById('rental-history').innerHTML = '<p class="error">An error occurred while fetching rental history.</p>';
        });
}

// Function to display rental history on the page
function displayRentalHistory(rentals) {
    const rentalHistoryDiv = document.getElementById('rental-history');
    rentalHistoryDiv.innerHTML = ''; // Clear any previous content

    rentals.forEach(rental => {
        const rentalDiv = document.createElement('div');
        rentalDiv.classList.add('booking');

        rentalDiv.innerHTML = `
            <h3>Booking ID: ${rental.booking_id}</h3>
            <p><strong>Vehicle ID:</strong> ${rental.vehicle_id}</p>
            <p><strong>Start Time:</strong> ${new Date(rental.start_time).toLocaleString()}</p>
            <p><strong>End Time:</strong> ${new Date(rental.end_time).toLocaleString()}</p>
            <p><strong>Status:</strong> ${rental.status}</p>
            <p><strong>Total Cost:</strong> $${parseFloat(rental.total_cost).toFixed(2)}</p>
            <p><strong>Created At:</strong> ${new Date(rental.created_at).toLocaleString()}</p>
            <p><strong>Updated At:</strong> ${new Date(rental.updated_at).toLocaleString()}</p>
        `;

        rentalHistoryDiv.appendChild(rentalDiv);
    });
}
