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
