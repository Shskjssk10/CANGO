# Payment Microservice Documentation

## Introduction

This API handles creation, retrieval and management of payment records. Stripe API is integrated for safer and more secure payment processing. 

## Base URL
```
http://localhost:8002/api/v1/
```

## API Endpoints

The following is the list of endpoints for this API:
- [Test DB Connection](#test-db-connection)
- [Post Payment](#post-payment)
- [Send Receipt](#send-receipt)
- [Create Payment Intent](#create-payment-intent)

---

### Test DB Connection
- Endpoints: `/test`
- Method: `GET`
- Description: Checks whether database connection is successful
- Response:
    - `200 OK`: Database connected successfully
    - `500 Internal Server Error`: Failed to connect to database

---

### Post Payment

- Endpoints: `/payment`
- Method: `POST`
- Description: Post payment for a booking
- Request Body:
```json
{
    "Amount": 100,
    "UserID": 4,
    "CarID": 1
}
```
- Response:
    - `200 OK`: Returns successful message.
    - `500 Internal Server Error`: Failed to post payment

---

### Send Receipt 
- Endpoints: `/paymentConfirmation`
- Method: `POST`
- Description: Send payment receipt code to user
-  Request Body:
```json
{
    "Name": "John Doe",
    "EmailAddr": "johndoe@example.com",
    "Model": "Model X",
    "Date": "2024-06-01",
    "StartTime": "09:00:00",
    "EndTime": "18:00:00",
    "Amount": 200
}
```
- Response:
    - `200 OK`: Receipt sent successfully

---

### Create Payment Intent
- Endpoints: `/create-payment-intent`
- Method: `POST`
- Description: Send payment intent to stripe API
-  Request Body:
```json
{
    "Amount": "1000",
    "Currency": "SGD"
}
```
- Response:
    - `200 OK`: Payment created successfully
    - `500 Internal Server Error`: Failed to create payment intent
---

## Return

[Click on this link to return to parent README.md](../../README.md)