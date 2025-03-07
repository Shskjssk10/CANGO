# CANGO

> [!NOTE]
> This repository is the local version of CANGO application! The deployed version is underway!

## Introduction 
In an era marked by sustainable transportation and shared economies, electric carsharing platforms have emerged as a cornerstone of modern urban mobility. This project aims to design and implement a fully functional electric car-sharing system using Go, with features catering to diverse user needs and real-world application scenarios. With an emphasis on practical and scalable solutions, the system includes user membership tiers, promotional discounts, and an accurate billing mechanism. 

## Table of Contents
1. [Architecture Diagram](#architecture-diagram)
2. [Design Considerations](#design-considerations)
    - [Microservice Breakdown](#1-microservice-breakdown)
    - [Database Design](#2-shared-database)
    - [Security](#3-security)
3. [Installation Guide](#installation-guide)
4. [API Documentation](#api-documentation)
5. [Future Implementations](#future-implementations)

## Architecture Diagram

 ![Updated-Architecture-Diagram](/client/images/architecture-diagram.png)

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


### **2. Shared Database**

All microservices share a single database, where all read and write operations are done. This is done for the following reasons:

* **Data-Consistency** 
    * All microservices operate on the same data
    * Minimises risk of data inconsistencies and conflicts 

* **Simplified Data Access** 
    * Simplifies data access logic for microservices
    * Minimises risk of data inconsistencies and conflicts
    * Eliminates the complexity of implementing data synchronization mechanisms between microservices.

### **3. Security**
* **Authentication:** Passwords are hashed before stored into database. This ensures that even in the scenario of a security breach, malicious users are unable to sign in as other users. Password hashing are done using bcrypt. 
* **Verification:** Verification is conducted when a user first signs up for an account, where a verification code is sent to the user's email. This code is then hashed when inputted and compared to the hashed code in the system. This ensures that the user is verified and is the legitimate owner signing up for an account, adding an additional layer of security.

## Installation Guide

1. Clone Repository 
```bash 
git clone https://github.com/Shskjssk10/CANGO.git
```
2. Copy and paste database code into MySQL Workbench
3. Add .env file under the CNAD_Assg1/server folder with the necessary keys (Please contact me. Unless you are my teacher grading it, it is together with the Github link). Your .env file should look something like below: 
```env
EMAIL_KEY = my_email_key
STRIPE_KEY = my_stripe_key
DB_USER = enter_db_user*
DB_PASS = enter_db_password*
DB_NAME = enter_db_name*

# All with '*' are to be inputted on your own
```
4. Navigate to the following directory 
```
cd \CNAD_Assg1\
```
5. Run the following command 
```bash
Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy Unrestricted
```
6. Run 
```bash
.\start-services.ps1
```
7. Navigate to client-side and run index.html with live server. Alternatively, enter  the following command:
```bash
cd .\client\
```
8. To kill the microservices, kill the terminal

> [!NOTE]
> Ensure Moesif Origin/CORS Changer extension is activated!

## API Documentation

Detailed API documentation for each microservice is linked below: 

- [Authentication Service API Documentation](./server/auth-service/auth_service_doc.md)
    - Handles user login, authentication, verification and registering of new accounts.
- [Payment Service API Documentation](./server/payment-service/payment_service_doc.md)
    - Handles payments, send receipts, Stripe API integrated for secure and seamless payment
- [User Management Service API Documentation](./server/user-management-service/user_management_service_doc.md)
    - Retrieves user info, updates new info
- [Vehicle Service API Documentation](./server/vehicle-service/vehicle_service_doc.md)
    - Retrieves car info, booking, update of booking



## Future Implementations

Proposed AWS Architecture Diagram 

 ![AWS-Architecture-Diagram](/client/images/aws-architecture-diagram.png)

- Host database on cloud
- Containerise application and deploy on cloud