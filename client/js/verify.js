// verify.js
document.getElementById("send-code").addEventListener("click", async () => {
    alert("Verification code has been sent to your email!");
    setTimeout(() => {
        console.log("6-digit code sent to user's email.");
    }, 1000);

    // Preparing Data
    const postData = {
        "Email": sessionStorage.getItem("EmailAddr")
    }; 

    // Verification Code Sent
    try {
        const response = await fetch("http://127.0.0.1:8000/api/v1/sendVerificationEmail", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(postData),
        });

        if (!response.ok) {
            throw new Error(`Failed to send verification email`);
        }

        const result = await response.json(); // Parse the JSON response
        console.log("Email has been sent successfully:", result);
    } catch (error) {
        console.error("Error:", error.message);
        alert("Registration Failed")
    }
});

document.getElementById("verify-form").addEventListener("submit", async (event) => {
    event.preventDefault();
    const code = document.getElementById("code").value;

    // Verify Verification code
    if (code.length === 6 && !isNaN(code)) {
        // Preparing Data
        const postData = {
            "Email": sessionStorage.getItem("EmailAddr"),
            "VerificationCode": code
        }; 

        try {
            const response = await fetch(`http://127.0.0.1:8000/api/v1/activateAccount`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(postData),
            });
    
            if (!response.ok) {
                alert("Verification Unsuccessful");
                throw new Error(`Verification Failed`);
            }
    
            const result = await response.json(); // Parse the JSON response
            console.log("Verification successfully:", result);
            window.location.href="verification_sucess.html"
        } catch (error) {
            console.error("Error:", error.message);
        }
    } else {
        alert("Please enter a valid 6-digit verification code.");
    }
});
