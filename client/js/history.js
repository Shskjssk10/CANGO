document.addEventListener('DOMContentLoaded', async () => {
    const historyCardsContainer = document.getElementById('history-cards');

    // Fetch booking history for the user
    const user = JSON.parse(sessionStorage.getItem("User"));
    try {
        const response = await fetch(`http://127.0.0.1:8001/api/v1/booking/user/${user.UserID}`);
        if (!response.ok) {
            throw new Error(`Error fetching bookings: ${response.statusText}`);
        }

        const bookings = await response.json();
        bookings.forEach(booking => {

            const card = createBookingCard(booking);
            historyCardsContainer.appendChild(card);
        });
    } catch (error) {
        console.error("Error:", error.message);
        historyCardsContainer.innerHTML = `<p>Unable to load booking history. Please try again later.</p>`;
    }

    // Create a booking card
    function createBookingCard(booking) {
        const card = document.createElement('div');
        card.classList.add('card');

        card.innerHTML = `
            <img src="../images/car.png" alt="Car Model" class="card-image">
            <div class="card-details">
                <h2 class="card-title">${booking.Model}</h2>
                <p class="card-time">Date ${booking.Date}</p>
                <p class="card-date">Start Time: ${booking.StartTime}</p>
                <p class="card-time">End Time: ${booking.EndTime}</p>
                <div class="card-buttons">
                    <button class="update-btn" onclick="updateBooking('${booking.BookingID}')">Update</button>
                    <button class="delete-btn" onclick="deleteBooking('${booking.BookingID}')">Delete</button>
                </div>
            </div>
        `;
        return card;
    }

    // Update booking
    window.updateBooking = (bookingID) => {
        sessionStorage.setItem("BookingID", bookingID);
        window.location.href="update-booking.html"
    };

    // Delete booking
    window.deleteBooking = async (bookingID) => {
        const confirmDelete = confirm("Are you sure you want to delete this booking?");
        if (confirmDelete) {
            try {
                const response = await fetch(`http://127.0.0.1:8001/api/v1/booking/${bookingID}`, {
                    method: 'DELETE',
                });

                if (!response.ok) {
                    throw new Error(`Error deleting booking: ${response.statusText}`);
                }

                alert("Booking deleted successfully!");
                location.reload(); // Reload the page to reflect changes
            } catch (error) {
                console.error("Error:", error.message);
                alert("Unable to delete booking. Please try again later.");
            }
        }
    };
});
