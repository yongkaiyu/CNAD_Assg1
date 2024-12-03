document.addEventListener('DOMContentLoaded', () => {

    // Call fetchVehicles when DOM content is fully loaded
    fetchVehicles();

});

// Fetch available vehicles and display them
async function fetchVehicles() {
    try {
        
        const response = await fetch('/api/v1/booking/vehicles');
        const vehicles = await response.json();

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
                <strong>License Plate:</strong> ${vehicle.licensePlate}<br>
                <strong>Location:</strong> ${vehicle.location}<br>
                <strong>Charge Level:</strong> ${vehicle.chargeLevel}%<br>
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

            bookButton.onclick = () => redirectToBooking(vehicle.id);
        });
    } catch (error) {
        console.error('Error fetching vehicles:', error);
    }
}

// Redirect to booking page with vehicle ID
function redirectToBooking(vehicleId) {
    window.location.href = `../booking?vehicleId=${vehicleId}`;
}