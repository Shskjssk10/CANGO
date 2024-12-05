document.addEventListener('DOMContentLoaded', async () => {
    const timeSlots = document.querySelectorAll('.time-slots input');
    const submitButton = document.querySelector('.submit-btn');
    const carPlaceholder = document.querySelector("#car-placeholder");
    const amountPayable = document.getElementById('amount-payable');
    const membershipTierDisplay = document.querySelector("#membership-tier")

    const carID = parseInt(sessionStorage.getItem("CarID"));
    let car;

    // Fetch car details
    try {
        const response = await fetch(`http://127.0.0.1:8001/api/v1/car/${carID}`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json"
            }
        });
        if (!response.ok) {
            throw new Error(`Network response was not ok: ${response.status} ${response.statusText}`);
        }
        car = await response.json();
        console.log(car);
    } catch (error) {
        console.log("Error:", error.message);
    }

    // Update placeholder with car model
    carPlaceholder.innerHTML = `Book ${car.Model}`;

    const user = JSON.parse(sessionStorage.getItem("User"));
    if (user.MembershipTier == "Basic") {
        membershipTierDisplay.innerHTML = `Basics (<p id="discount">No</p> Discount)`
    } else if (user.MembershipTier == "Premium") {
        membershipTierDisplay.innerHTML = `Premium (<p id="discount">5</p>% Discount)`
    } else {
        membershipTierDisplay.innerHTML = `VIP (<p id="discount">8</p>% Discount)`
    }

    const discount = document.querySelector("#discount").innerHTML;

    // Set car JSON for future use
    sessionStorage.setItem("Car", JSON.stringify(car));

    // Add event listener to time slots for price calculation
    timeSlots.forEach(slot => {
        slot.addEventListener('change', () => {
            updatePrice();
        });
    });

    // Event listener for booking confirmation
    submitButton.addEventListener('click', () => {
        const selectedSlots = [];
        timeSlots.forEach(slot => {
            if (slot.checked) {
                selectedSlots.push(slot.value);
            }
        });
        const amount = parseInt(document.querySelector("#amount-payable").innerHTML);
        console.log(amount)
        const date = document.getElementById('date').value;

        // Ensure all selected time slots are consecutive
        if (!date) {
            alert('Please select a date for the booking.');
            return;
        }
        if (selectedSlots.length > 0 && validateConsecutiveSlots(selectedSlots)) {
            const startTime = selectedSlots[0];
            const endTime = calculateEndTime(selectedSlots);

            alert(`Booking is valid for ${date} from ${startTime} to ${endTime}. Total payable: $${calculateTotalPrice(selectedSlots.length)}. Redirecting to payment`);

            const booking = {
                "Date": date,
                "StartTime": startTime,
                "EndTime": endTime,
                "UserID": user.UserID,
                "CarID": parseInt(carID),
                "PaymentID": null
            };

            const payment = {
                "Amount": parseInt(amount),
                "UserID": user.UserID,
                "CarID": parseInt(carID)
            }

            sessionStorage.setItem("BookingDetails", JSON.stringify(booking))
            sessionStorage.setItem("PaymentDetails", JSON.stringify(payment))

            // console.log("Booking Details:", booking);
            // console.log("Payment Details:", payment);

            window.location.href="payment.html"
            // Perform further actions here, like sending the booking to the server
        } else if (selectedSlots.length === 0) {
            alert('Please select at least one time slot.');
        } else {
            alert('Selected time slots must be consecutive.');
        }
    });

    // Update price based on selected slots
    function updatePrice() {
        const selectedCount = Array.from(timeSlots).filter(slot => slot.checked).length;
        const totalPrice = calculateTotalPrice(selectedCount);
        amountPayable.textContent = `${totalPrice}`;
    }

    // Calculate total price
    function calculateTotalPrice(slotCount) {
        const ratePerSession = car.RentalRate;
        if (discount == "No") {
            return slotCount * ratePerSession;
        } else {
            return slotCount * ratePerSession * (parseInt(discount)/100)
        }
    }

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
