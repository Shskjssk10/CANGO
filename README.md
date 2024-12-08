# CNAD_Assg1

## Design Considerations

**1. Microservice Breakdown**
All microservices are stored under the 'server' folder

* **auth-service:**
    * Handles user registration, login, and authentication.

* **payment-service:**
    * Handles the creation, retrieval, and management of payment records.
    * Interacts with the database to store and retrieve payment information.
    * May integrate with other services for order processing or notifications.

* **stripe-service:**
    * Responsible for handling secure payment transactions using the Stripe API.
    * Integrates with Stripe for payment processing, refunds, and other payment-related operations.
    * Interacts with the payment-service to record successful or failed payment transactions.

* **user-management-service:**
    * Handles user profile updates, including name, contact information, etc.
    * Retrieves user-specific information based on user email address.

* **vehicle-service:**
    * Handles fetching and managing car and booking data.
    * Provides APIs for retrieving vehicle information, searching for available vehicles, and managing bookings.
    * Interacts with the payment-service for user payment of booking.


**2. Shared Database:**
* All microservices share a single database for data consistency and efficient data access.

**3. Independent Deployment:**
* Each microservice is deployed independently, allowing for flexible scaling and updates without affecting other services.

**4. Security:**
* Passwords are hashed before stored into database. Sending verification codes via email is also implemented to ensure that the use of email
address is limited to the owner of the email. 

