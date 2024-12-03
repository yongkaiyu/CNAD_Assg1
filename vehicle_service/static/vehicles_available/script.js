document.addEventListener('DOMContentLoaded', () => {

    // Call fetchVehicles when DOM content is fully loaded
    fetchVehicles();

});

// Fetch available vehicles and display them
async function fetchVehicles() {
    try {
        
        const response = await fetch('/api/v1/booking/vehicles', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            }
        });
        
        // Check if the response is successful
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Error fetching vehicles');
        }

        const vehicles = await response.json();

        // Debug: Check what we're receiving as vehicles
        console.log('Received vehicles:', vehicles);

        const vehicleList = document.getElementById('vehicle-list');
        vehicleList.innerHTML = ''; // Clear existing list

        if (vehicles.length === 0) {
            // Display a message when no vehicles are available
            const noVehiclesMessage = document.createElement('div');
            noVehiclesMessage.className = 'no-vehicles-message';
            noVehiclesMessage.textContent = 'No vehicles available for booking.';
            vehicleList.appendChild(noVehiclesMessage);
            return;
        }
        
        vehicles.forEach(vehicle => {
            // Create container for each vehicle
            const vehicleContainer = document.createElement('div');
            vehicleContainer.className = 'vehicle-container';

            // Add vehicle details
            const details = document.createElement('div');
            details.className = 'vehicle-details';
            details.innerHTML = `
                <strong>License Plate:</strong> ${vehicle.license_plate}<br>
                <strong>Location:</strong> ${vehicle.location}<br>
                <strong>Charge Level:</strong> ${vehicle.charge_level}%<br>
                <strong>Status:</strong> ${vehicle.status}<br>
            `;
            vehicleContainer.appendChild(details);

            // Add book button
            const bookButton = document.createElement('button');
            bookButton.className = 'book-button';
            bookButton.textContent = 'Book';
            vehicleContainer.appendChild(bookButton);

            // Append container to the list
            vehicleList.appendChild(vehicleContainer);

            bookButton.onclick = () => redirectToBooking(vehicle);
        });
    } catch (error) {
        console.error('Error fetching vehicles:', error);
        // Handle the error gracefully on the frontend
        const vehicleList = document.getElementById('vehicle-list');
        vehicleList.innerHTML = `<p>Error: ${error.message}</p>`;
    }
}

// Redirect to booking page with vehicle ID
function redirectToBooking(vehicle) {
    window.location.href = `../vehicle_booking/`;
    localStorage.setItem("vehicleId",vehicle.vehicle_id)
}

/*function redirectToBooking(vehicleId) {
    window.location.href = `../booking?vehicleId=${vehicleId}`;
}*/