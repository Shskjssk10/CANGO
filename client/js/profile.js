document.addEventListener('DOMContentLoaded', async () => {
    // Get User
    const email = sessionStorage.getItem("EmailAddr")
    let user;

    try {
        const response = await fetch(`http://127.0.0.1:8004/api/v1/getUser/${email}`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json"
            }
        }); // Adjust the URL as needed
        if (!response.ok) {
            throw new Error(`Network response was not ok: ${response.status} ${response.statusText}`);
        }
        user = await response.json();
        console.log(user);
    } catch (error) {
        console.log("Error:", error.message);
    }

    sessionStorage.setItem("User", JSON.stringify(user))


    const profileForm = document.getElementById('profileForm');
    const emailInput = document.getElementById('email');
    const contactNoInput = document.getElementById('contactno');
    const passwordInput = document.getElementById('password');
    const confirmPasswordInput = document.getElementById('confirmPassword');

    const membershipTier = document.querySelector("#membership-tier");

    membershipTier.innerHTML = `<strong>Membership Tier:</strong> ${user.MembershipTier}`

    emailInput.value = user.EmailAddr;
    contactNoInput.value = user.ContactNo;
    passwordInput.value = "";

    // Handle form submission
    profileForm.addEventListener('submit', async function (event) {
        event.preventDefault(); // Prevent default form submission

        const email = emailInput.value;
        const contactNo = contactNoInput.value;
        const password = passwordInput.value;
        const confirmPassword = confirmPasswordInput.value;

        // Validate inputs
        if (!email || !contactNo ) {
            alert('Please fill out all fields.');
            return;
        }

        if (!/^\d{8}$/.test(contactNo)) {
            alert('Please enter a valid 8-digit contact number.');
            return;
        }

        if ((password && confirmPassword) && (password !== confirmPassword)) {
            alert('Passwords do not match!');
            return;
        }

        // Here you would typically send the data to your API for updating

        const updateData = {
            "Name": user.Name,
            "ContactNo": contactNo,
            "EmailAddr": email,
            "PasswordHash": password
        }

        try {
            const response = await fetch(`http://127.0.0.1:8004/api/v1/update/${user.UserID}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(updateData),
            });

            // Check if the request was successful
            if (response.ok) {
                const result = await response.json();
                console.log('Profile updated successfully:', result);
                alert('Your profile has been updated!');
            } else {
                throw new Error('Failed to update profile.');
            }
        } catch (error) {
            console.error('Error:', error.message);
            alert('There was an error updating your profile. Please try again later.');
        }

        sessionStorage.setItem("EmailAddr", email)
        location.reload();
    });
});
