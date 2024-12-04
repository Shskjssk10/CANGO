// script.js
document.getElementById('registerForm').addEventListener('submit', async function(event) {
    event.preventDefault();
    
    const name = document.getElementById('name').value.trim();
    const email = document.getElementById('email').value.trim();
    const contactNo = document.getElementById('contactNo').value.trim();
    const password = document.getElementById('password').value.trim();
    const confirmPassword = document.getElementById('confirmPassword').value.trim();
    
    if (!name || !email || !contactNo || !password || !confirmPassword) {
        alert('All fields must be filled out.');
        return;
    }
    
    if (!/^\d{8}$/.test(contactNo)) {
        alert('Contact number must be an 8-digit number.');
        return;
    }
    
    if (password !== confirmPassword) {
        alert('Passwords do not match.');
        return;
    }

    // Preparing Data
    const userData = {
        "Name": name,
        "EmailAddr": email,
        "ContactNo": contactNo,
        "PasswordHash": confirmPassword
    };

    console.log(JSON.stringify(userData.name))
    sessionStorage.setItem("EmailAddr", userData.EmailAddr);
    
    // Register User
    try {
        const response = await fetch("http://127.0.0.1:8000/api/v1/registerUser", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(userData),
        });

        if (!response.ok) {
            throw new Error(`Failed to create user: ${response.status} ${response.statusText}`);
        }

        const result = await response.json(); // Parse the JSON response
        console.log("User created successfully:", result);
        window.location.href = "verify.html"
    } catch (error) {
        console.error("Error:", error.message);
        alert("Registration Failed")
    }
});
