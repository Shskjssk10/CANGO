document.addEventListener('DOMContentLoaded', async () => {
    const carList = document.getElementById('carList');

    // Getting Car Data
    let listOfCars;

    try {
        const response = await fetch(`http://127.0.0.1:8001/api/v1/cars`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json"
            }
        }); // Adjust the URL as needed
        if (!response.ok) {
            throw new Error(`Network response was not ok: ${response.status} ${response.statusText}`);
        }
        listOfCars = await response.json();
        console.log(listOfCars);
    } catch (error) {
        console.log("Error:", error.message);
    }

    // Generate car cards dynamically
    listOfCars.forEach(car => {
        const card = document.createElement('div');
        card.className = 'car-card';

        card.innerHTML = `
            <img src="../images/car.png" alt="${car.Model}">
            <div class="car-card-content">
                <h2>${car.Model}</h2>
                <p>CANGO newly bought car! Be sure to book it soon!!! </p>
                <p><strong>Rental Rate</strong>: $${car.RentalRate}/hr</p>
                <button class="view-details-btn">Book Now!</button>
            </div>
        `;

        carList.appendChild(card);
    });
});
