window.addEventListener("DOMContentLoaded", () => {
    let userId = localStorage.getItem("userId");

    // Trim and sanitize userId to remove unwanted characters
    userId = userId.trim().replace(/[^0-9]/g, "");

    if (!userId) {
        alert("Invalid User ID provided.");
        return;
    }

    fetch(`http://localhost:5000/api/v1/billing/bills?user_id=${userId}`)
        .then((response) => {
            if (!response.ok) {
                throw new Error(`Failed to fetch billing information: ${response.statusText}`);
            }
            return response.json();
        })
        .then((data) => {
            console.log("Billing data:", data); // Add this to debug the response
            const tbody = document.querySelector("#billing-table tbody");
            tbody.innerHTML = ""; // Clear existing rows

            if (!Array.isArray(data) || data.length === 0) {
                tbody.innerHTML = `<tr><td colspan="7">No billing records found.</td></tr>`;
                return;
            }

            data.forEach((record) => {
                const row = document.createElement("tr");
                row.innerHTML = `
                    <td>${record.billing_id}</td>
                    <td>${record.booking_id}</td>
                    <td>${record.payment_status}</td>
                    <td>${record.payment_method}</td>
                    <td>${record.total_amount.toFixed(2)}</td>
                    <td>${new Date(record.created_at).toLocaleString()}</td>
                    <td>${new Date(record.updated_at).toLocaleString()}</td>
                `;
                tbody.appendChild(row);
            });
        })
        .catch((error) => {
            console.error("Error:", error);
            alert("An error occurred while fetching billing information.");
        });
});
