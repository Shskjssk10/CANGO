# User Management Microservice Documentation

## Introduction

This API handles user info retrieval and user profile updates such as passwords and contact numbers.

## Base URL
```
http://localhost:8004/api/v1/
```

## API Endpoints

The following is the list of endpoints for this API:
- [Test DB Connection](#test-db-connection)
- [Get User Info](#get-user-info)
- [Update User Info](#update-user-info)

---
### Test DB Connection
- Endpoints: `/test`
- Method: `GET`
- Description: Checks whether database connection is successful
- Response:
    - `200 OK`: Database connected successfully
    - `500 Internal Server Error`: Failed to connect to database

---

### Get User Info
- Endpoints: `/getUser/{email}`
- Method: `GET`
- Description: Gets user information
- Path Parameter: 
    - `email`: The unique email entered by each user.
- Response:
    - `200 OK`: Returns user information in JSON format.
    ```json
        {
            "UserID": 1,
            "Name": "John Doe",
            "EmailAddr": "johndoe@gmail.com",
            "ContactNo": "21212121",
            "MembershipTier": "Premium",
            "PasswordHash": "hashed_password1",
            "IsActivated": 1,
            "VerificationCodeHash": "hash2"
        }
    ```
    - `500 Internal Server Error`: Failed to retrieve user information

---

### Update User Info
- Endpoints: `/update/{id}`
- Method: `POST`
- Description: Update information of each user
- Path Parameter: 
    - `id`: The unique identifier of each user.
- Request Body:
```json
{
    "Name": "John Doe",
    "ContactNo": "11112222",
    "EmailAddr": "newemail@example.com",
    "PasswordHash": "newpassword123!"
}
```
- Response:
    - `200 OK`: User has been updated.
    - `500 Internal Server Error`: Email or contact no. already exists

---

## Return

[Click on this link to return to parent README.md](../../README.md)