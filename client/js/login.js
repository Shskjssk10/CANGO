// script.js
document.addEventListener('DOMContentLoaded', async () => {
    document.getElementById('loginForm').addEventListener('submit', async function(event) {
        event.preventDefault();
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;

        const loginData = {
            "Email": email,
            "Password": password
        }

        if (email && password) {
            try {
                const response = await fetch("http://127.0.0.1:8000/api/v1/loginUser", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(loginData),
                });
        
                if (!response.ok) {
                    throw new Error(`Failed to login: ${response.status} ${response.statusText}`);
                }
        
                const result = await response.json(); // Parse the JSON response
                console.log("Login successfully:", result);
                alert("Login Successful")
                window.location.href = "dashboard.html"
            } catch (error) {
                alert("Login Unsuccessful - Wrong Email or Password")
                console.error("Error:", error.message);
            }
        } else {
            alert('Please fill out all fields.');
        }
    });

    document.getElementById('registerBtn').addEventListener('click', function() {
        alert('Redirecting to the registration page.');
        // Simulate redirect
        window.location.href = '/register'; // Replace with actual registration page URL
    });
});
