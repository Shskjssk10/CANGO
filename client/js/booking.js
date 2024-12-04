document.addEventListener('DOMContentLoaded', async () => {
    const timeSlots = document.querySelectorAll('.time-slots input');
    const submitButton = document.querySelector('.submit-btn');
    
    const carPlaceholder = document.querySelector("#car-placeholder");

    const carID = sessionStorage.getItem("CarID");

    // Get Car Details
    let car;
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

    // Edit placeholder
    carPlaceholder.innerHTML = `Book ${car.Model}`

    // Set Car JSON for Payment
    sessionStorage.setItem("Car", JSON.stringify(car));

    submitButton.addEventListener('click', () => {
        const selectedSlots = [];
        timeSlots.forEach(slot => {
            if (slot.checked) {
                selectedSlots.push(slot.value);
            }
        });

        // Ensure all selected time slots are consecutive
        if (selectedSlots.length > 0 && validateConsecutiveSlots(selectedSlots)) {
            alert(`Booking confirmed for the following slots: ${selectedSlots.join(', ')}`);
            // Redirect to confirmation or perform additional actions here
        } else if (selectedSlots.length === 0) {
            alert('Please select at least one time slot.');
        } else {
            alert('Selected time slots must be consecutive.');
        }
    });

    // Function to validate that selected slots are consecutive
    function validateConsecutiveSlots(slots) {
        const timeOrder = ['0900-1200', '1200-1500', '1500-1800', '1800-2100'];
        const slotIndexes = slots.map(slot => timeOrder.indexOf(slot));

        // Check if the slots form a consecutive sequence
        for (let i = 0; i < slotIndexes.length - 1; i++) {
            if (slotIndexes[i + 1] !== slotIndexes[i] + 1) {
                return false;
            }
        }
        return true;
    }
});
