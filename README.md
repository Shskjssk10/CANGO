# CANGO
In an era marked by sustainable transportation and shared economies, electric carsharing platforms have emerged as a cornerstone of modern urban mobility. This project aims to design and implement a fully functional electric car-sharing system using Go, with features catering to diverse user needs and real-world application scenarios. With an emphasis on practical and scalable solutions, the system includes user membership tiers, promotional discounts, and an accurate billing mechanism. 

## Architecture Diagram

 ![Updated-Architecture-Diagram](/client/images/Updated-Architecture-Diagram.png)

## Design Considerations

### **1. Microservice Breakdown**

All microservices are stored under the 'server' folder

* **auth-service:**
    * Handles user registration, login, and authentication.

* **payment-service:**
    * Handles the creation, retrieval, and management of payment records.
    * Interacts with the database to store and retrieve payment information.
    * May integrate with other services for order processing or notifications.
    * Integrates Stripe API for handling secure payment transactions
    * Integrates with Stripe for payment processing, refunds, and other payment-related operations.

* **user-management-service:**
    * Handles user profile updates, including name, contact information, etc.
    * Retrieves user-specific information based on user email address.

* **vehicle-service:**
    * Handles fetching and managing car and booking data.
    * Provides APIs for retrieving vehicle information, searching for available vehicles, and managing bookings.
    * Interacts with the payment-service for user payment of booking.


**2. Shared Database:**

All microservices share a single database, where all read and write operations are done. This is done for the following reasons:

* **Data-Consistency** 
    * All microservices operate on the same data
    * Minimises risk of data inconsistencies and conflicts 

* **Simplified Data Access** 
    * Simplifies data access logic for microservices
    * Minimises risk of data inconsistencies and conflicts
    * Eliminates the complexity of implementing data synchronization mechanisms between microservices.

**3. Independent Deployment:**
* Each microservice is deployed independently, allowing for flexible scaling and updates without affecting other services.

**4. Security:**
* **Authentication:** Passwords are hashed before stored into database. This ensures that even in the scenario of a security breach, malicious users are unable to sign in as other users. Password hashing are done using bcrypt. 
* **Verification:** Verification is conducted when a user first signs up for an account, where a verification code is sent to the user's email. This code is then hashed when inputted and compared to the hashed code in the system. This ensures that the user is verified and is the legitimate owner signing up for an account, adding an additional layer of security.

## Instructions for setting up

1. Clone Repository
2. Run database in MySQL Workbench
3. Add .env file under the CNAD_Assg1/server folder with the necessary keys (Please contact me. Unless you are my teacher grading it, it is together with the Github link)
4. Navigate to the following directory 

```cd .\CNAD_Assg1\```
5. Run the following command 
```Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy Unrestricted```
6. Run 
```.\start-services.ps1```
7. Navigate to ```cd .\client\html\``` and run index.html with live server.
8. To kill the microservices, kill the terminal

> [!NOTE]
> Ensure Moesif Origin/CORS Changer extension is activated!

## Future Implementations

- Host database on cloud
- Containerise application and deploy on cloud