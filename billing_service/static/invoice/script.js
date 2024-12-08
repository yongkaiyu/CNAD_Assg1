document.addEventListener('DOMContentLoaded', function () {
    const bookingId = localStorage.getItem('bookingInvoiceId');

    if (!bookingId) {
        document.getElementById('invoice-details').innerHTML = '<p class="error">No booking ID provided.</p>';
        return;
    }

    fetch(`/api/v1/billing/invoice?booking_id=${bookingId}`)
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                document.getElementById('invoice-details').innerHTML = `<p class="error">${data.error}</p>`;
            } else {
                displayInvoice(data);
            }
        })
        .catch(error => {
            console.error('Error fetching invoice:', error);
            document.getElementById('invoice-details').innerHTML = '<p class="error">An error occurred while fetching invoice details.</p>';
        });
});

function displayInvoice(invoice) {
    const invoiceDiv = document.getElementById('invoice-details');
    invoiceDiv.innerHTML = `
        <h1>Invoice</h1>
        <p><strong>Booking ID:</strong> ${invoice.booking_id}</p>
        <p><strong>User ID:</strong> ${invoice.user_id}</p>
        <p><strong>Vehicle ID:</strong> ${invoice.vehicle_id}</p>
        <p><strong>Start Time:</strong> ${new Date(invoice.start_time).toLocaleString()}</p>
        <p><strong>End Time:</strong> ${new Date(invoice.end_time).toLocaleString()}</p>
        <p><strong>Total Cost:</strong> $${parseFloat(invoice.total_cost).toFixed(2)}</p>
        <p><strong>Status:</strong> ${invoice.status}</p>
        <p><strong>Created At:</strong> ${new Date(invoice.created_at).toLocaleString()}</p>
        <p><strong>Updated At:</strong> ${new Date(invoice.updated_at).toLocaleString()}</p>
        <p><strong>Generated At:</strong> ${new Date(invoice.generated_at).toLocaleString()}</p>
    `;
}
