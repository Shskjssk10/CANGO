# Authentication Microservice Documentation

## Introduction

This API handles user registration, user login, user authentication, sending of verification codes and verifying said verification codes.

## Base URL
```
http://localhost:8000/api/v1/
```

## API Endpoints

The following is the list of endpoints for this API:
- [Test DB Connection](#test-db-connection)
- [Register User](#register-user)
- [Login User](#login-user)
- [Send Verification Code](#send-verification-email)
- [Verify Verification Code](#verify-verification-code)

---
### Test DB Connection
- Endpoints: `/test`
- Method: `GET`
- Description: Checks whether database connection is successful
- Response:
    - `200 OK`: Database connected successfully
    - `500 Internal Server Error`: Failed to connect to database

---

### Register User
- Endpoints: `/registerUser`
- Method: `POST`
- Description: Registers a new user
- Request Body:
```json
{
	"Name": "John Doe",
	"EmailAddr":"user@example.com",
	"ContactNo": "10008000",
	"PasswordHash": "password123!"
}
```
- Response:
    - `200 OK`: Database connected successfully
    - `400 Bad Request`: Invalid Request Body
    - `500 Internal Server Error`: Email or contact no. already exists

---

### Login User
- Endpoints: `/loginUser`
- Method: `POST`
- Description: Logins a user
-  Request Body:
```json
{
    "Email": "johndoe@example.com",
    "Password": "password123!"
}
```
- Response:
    - `200 OK`: Login successful
    - `401 Unauthorised`: Invalid email or password

---

### Send Verification Email
- Endpoints: `/sendVerificationEmail`
- Method: `POST`
- Description: Send verification code to user
-  Request Body:
```json
{
    "Email": "johndoe@example.com"
}
```
- Response:
    - `200 OK`: Login successful

---

### Verify Verification Code
- Endpoints: `/activateAccount`
- Method: `PUT`
- Description: Verify verification code sent
-  Request Body:
```json
{
    "Email": "johndoe@example.com",
    "VerificationCode": "123456"
}
```
- Response:
    - `200 OK`: Verification Successful
    - `401 Unauthorised`: Invalid verification code

---