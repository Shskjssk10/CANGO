document.addEventListener('DOMContentLoaded', async () => {

    document.getElementById("verify-form").addEventListener("submit", (event) => {
        event.preventDefault();
        window.location.href = "homepage.html"
    });

    let booking = JSON.parse(sessionStorage.getItem("BookingDetails"));
    const payment = JSON.parse(sessionStorage.getItem("PaymentDetails"));

    console.log(payment)

    // Post Payment
    let paymentResult;
    try {
        const response = await fetch("http://127.0.0.1:8002/api/v1/payment", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                "Amount": payment.Amount,
                "CarID": payment.CarID,
                "UserID": payment.UserID
            }),
        });

        if (!response.ok) {
            throw new Error(`Failed to create payment: ${response.status} ${response.statusText}`);
        }

        paymentResult = await response.json(); // Parse the JSON response
        console.log("Payment created successfully:", paymentResult);
    } catch (error) {
        console.error("Error:", error.message);
        alert("Registration Failed")
    }

    booking.PaymentID = parseInt(paymentResult.PaymentID)
    console.log(booking);

    // Post Booking
    try {
        const response = await fetch("http://127.0.0.1:8001/api/v1/booking", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(booking),
        });

        if (!response.ok) {
            throw new Error(`Failed to create booking: ${response.status} ${response.statusText}`);
        }

        const result = await response.json(); // Parse the JSON response
        console.log("Booking created successfully:", result);
    } catch (error) {
        console.error("Error:", error.message);
        alert("Registration Failed")
    }

    // Send Confirmation Email
    const confirmationDetails = {
        "PaymentID": parseInt(paymentResult.PaymentID),
        "Amount": payment.Amount,
        "UserID": booking.UserID,
        "CarID": booking.CarID
    }

    console.log(confirmationDetails);
    try {
        const response = await fetch("http://127.0.0.1:8002/api/v1/paymentConfirmation", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(confirmationDetails),
        });

        if (!response.ok) {
            throw new Error(`Failed to send email: ${response.status} ${response.statusText}`);
        }

        const result = await response.json(); // Parse the JSON response
        console.log("Email sent successfully:", result);
    } catch (error) {
        console.error("Error:", error.message);
        alert("Registration Failed")
    }
});

