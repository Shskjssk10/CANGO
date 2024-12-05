document.addEventListener('DOMContentLoaded', async () => {
    console.log("Homepage loaded successfully!");

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

    document.querySelector("#welcome").innerHTML = `Welcome ${user.Name}`
    
});
