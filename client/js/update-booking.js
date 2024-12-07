document.addEventListener('DOMContentLoaded', async () => {
    const timeSlots = document.querySelectorAll('.time-slots input');
    const submitButton = document.querySelector('.submit-btn');
    const carPlaceholder = document.querySelector("#car-placeholder");

    const bookingID = parseInt(sessionStorage.getItem("BookingID"));
    let booking;

    // Fetch Booking details
    try {
        const response = await fetch(`http://127.0.0.1:8001/api/v1/booking/${bookingID}`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json"
            }
        });
        if (!response.ok) {
            throw new Error(`Network response was not ok: ${response.status} ${response.statusText}`);
        }
        booking = await response.json();
        console.log(booking);
    } catch (error) {
        console.log("Error:", error.message);
    }

    // Event listener for booking confirmation
    submitButton.addEventListener('click', async () => {
        const selectedSlots = [];
        timeSlots.forEach(slot => {
            if (slot.checked) {
                selectedSlots.push(slot.value);
            }
        });
        const date = document.getElementById('date').value;

        // Ensure all selected time slots are consecutive
        if (!date) {
            alert('Please select a date for the booking.');
            return;
        }
        if (selectedSlots.length > 0 && validateConsecutiveSlots(selectedSlots)) {
            const startTime = selectedSlots[0];
            const endTime = calculateEndTime(selectedSlots);

            const details = {
                "Date": date,
                "StartTime": startTime, 
                "EndTime": endTime,
                "CarID": booking.CarID
            }

            // Check Booking Validity
            try {
                const response = await fetch("http://127.0.0.1:8001/api/v1/checkValidity", {
                    method: "PUT",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(details),
                });

                if (!response.ok) {
                    throw new Error(`Failed to check booking: ${response.status} ${response.statusText}`);
                }

                const result = await response.json(); // Parse the JSON response
                console.log("Booking Checked successfully:", result);

                // If Valid
                if (result.StatusCode == 200) {
                    alert(`Booking is valid for ${date} from ${startTime} to ${endTime}. Booking is now being updated!`);

                    const newBooking = {
                        "StartTime": startTime,
                        "EndTime": endTime,
                        "Date": date
                    }
        
                    // Update Booking
                    try {
                        const response = await fetch(`http://127.0.0.1:8001/api/v1/booking/${bookingID}`, {
                            method: 'PUT',
                            headers: {
                                'Content-Type': 'application/json',
                            },
                            body: JSON.stringify(newBooking),
                        });
                        if (!response.ok) {
                            throw new Error(`Network response was not ok: ${response.status} ${response.statusText}`);
                        }
                    } catch (error) {
                        console.log("Error:", error.message);
                    }
        
                    alert("Booking has been successfully updated!")
                    window.location.href="view-history.html"
        
                } else {
                    alert("This booking is invalid as it clashes with another booking!")
                }

            } catch (error) {
                console.error("Error:", error.message);
                alert("Registration Failed")
            }
        } else if (selectedSlots.length === 0) {
            alert('Please select at least one time slot.');
        } else {
            alert('Selected time slots must be consecutive.');
        }
    });

    // Validate that selected slots are consecutive
    function validateConsecutiveSlots(slots) {
        const timeOrder = ['09:00:00', '12:00:00', '15:00:00', '18:00:00'];
        const slotIndexes = slots.map(slot => timeOrder.indexOf(slot));

        for (let i = 0; i < slotIndexes.length - 1; i++) {
            if (slotIndexes[i + 1] !== slotIndexes[i] + 1) {
                return false;
            }
        }
        return true;
    }

    // Calculate End Time based on the last selected time slot
    function calculateEndTime(slots) {
        const timeOrder = ['09:00:00', '12:00:00', '15:00:00', '18:00:00'];
        const lastSlot = slots[slots.length - 1];
        const lastIndex = timeOrder.indexOf(lastSlot);
        return lastIndex < timeOrder.length - 1 ? timeOrder[lastIndex + 1] : '21:00:00'; // End time for the last slot
    }
});
