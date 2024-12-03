document.addEventListener("DOMContentLoaded", async () => {
    const userId = localStorage.getItem("userId");
    if (!userId) {
        alert("User not logged in!");
        return;
    }

    try {
        // Fetch booked vehicles
        const response = await fetch(`/api/v1/booking/bookings?userId=${userId}`);
        const bookedVehicles = await response.json();

        const vehicleList = document.getElementById("vehicle-list");

        if (bookedVehicles.length === 0) {
            vehicleList.innerHTML = "<p>No booked vehicles found.</p>";
            return;
        }

        // Display each booked vehicle
        bookedVehicles.forEach(vehicle => {
            const vehicleContainer = document.createElement("div");
            vehicleContainer.className = "vehicle-container";

            const details = document.createElement("div");
            details.className = "vehicle-details";
            details.innerHTML = `
                <strong>License Plate:</strong> ${vehicle.licensePlate}<br>
                <strong>Location:</strong> ${vehicle.location}<br>
                <strong>Charge Level:</strong> ${vehicle.chargeLevel}%<br>
                <strong>Booking ID:</strong> ${vehicle.bookingId}<br>
                <strong>Start Time:</strong> ${vehicle.startTime}<br>
                <strong>End Time:</strong> ${vehicle.endTime}
            `;
            vehicleContainer.appendChild(details);

            const actions = document.createElement("div");
            actions.className = "actions";

            const modifyButton = document.createElement("button");
            modifyButton.className = "modify-button";
            modifyButton.textContent = "Modify Booking";
            modifyButton.onclick = () => redirectToModifyBooking(vehicle.bookingId);

            const deleteButton = document.createElement("button");
            deleteButton.className = "delete-button";
            deleteButton.textContent = "Delete Booking";
            deleteButton.onclick = () => deleteBooking(vehicle.bookingId);

            actions.appendChild(modifyButton);
            actions.appendChild(deleteButton);
            vehicleContainer.appendChild(actions);

            vehicleList.appendChild(vehicleContainer);
        });
    } catch (error) {
        console.error("Error fetching booked vehicles:", error);
        alert("An error occurred while loading booked vehicles.");
    }
});

// Redirect to Modify Booking Page
function redirectToModifyBooking(bookingId) {
    window.location.href = `/modify-booking.html?bookingId=${bookingId}`;
}

// Delete Booking
async function deleteBooking(bookingId) {
    if (confirm("Are you sure you want to cancel this booking?")) {
        try {
            const response = await fetch(`/api/v1/booking/cancel/${bookingId}`, {
                method: "DELETE",
            });

            if (response.ok) {
                alert("Booking deleted successfully.");
                window.location.reload(); // Reload the page to refresh the list
            } else {
                const error = await response.json();
                alert(`Error deleting booking: ${error.message}`);
            }
        } catch (error) {
            console.error("Error deleting booking:", error);
            alert("An error occurred while deleting the booking.");
        }
    }
}