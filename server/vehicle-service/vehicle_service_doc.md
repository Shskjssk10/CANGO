# Vehicle Microservice Documentation

## Introduction

This API handles fetching and managing of car and booking data

## Base URL
```
http://localhost:8001/api/v1/
```
## API Endpoints

The following is the list of endpoints for this API:
- [Test DB](#test-db-connection)
- [Get All Cars](#get-all-cars)
- [Get Car](#get-car)
- [Update Car Location](#update-car-location)
- [Post Booking](#post-booking)
- [Check Booking Validity](#check-booking-validity)
- [Get Booking](#get-booking)
- [Update Booking](#update-booking)
- [Delete Booking](#delete-booking)
- [Get Booking By Car ID](#get-booking-by-car-id)
- [Get Booking By User ID](#get-booking-by-user-id)

---
### Test DB Connection
- Endpoints: `/test`
- Method: `GET`
- Description: Checks whether database connection is successful
- Response:
    - `200 OK`: Database connected successfully
    - `500 Internal Server Error`: Failed to connect to database

--- 

### Get All Cars

- Endpoints: `/cars`
- Method: `GET`
- Description: Gets all car's information
- Response:
    - `200 OK`: Returns all cars information in JSON format.
    ```json
        [
            {
                "CarID": 1,
                "Model": "Toyota Camry",
                "PlateNo": "ABC1234",
                "RentalRate": 50,
                "Location": "Lorong Ah Soo"
            },
            {
                "CarID": 2,
                "Model": "Honda Civic",
                "PlateNo": "DEF5678",
                "RentalRate": 40,
                "Location": "Lorong Ah Soo"
            },(...and more)
        ]
    ```
    - `500 Internal Server Error`: Failed to retrieve cars' information

---

### Get Car

- Endpoints: `/car/{id}`
- Method: `GET`
- Description: Gets all car's information
- Path Parameter: 
    - `id`: The unique id for each car.
- Response:
    - `200 OK`: Returns all cars information in JSON format.
    ```json
        {
            "CarID": 3,
            "Model": "BMW Type X",
            "PlateNo": "DEF5678",
            "RentalRate": 40,
            "Location": "Lorong Ah Soo"
        }
    ```
    - `500 Internal Server Error`: Failed to retrieve car information

---

### Update Car Location

- Endpoints: `/car/{id}`
- Method: `PUT`
- Description: Update car's location
- Path Parameter: 
    - `id`: The unique id for each car.
- Request Body:
```json
{
    "Location": "Punggol"
}
```
- Response:
    - `200 OK`: Returns successful message.
    - `500 Internal Server Error`: Failed to update car information

---

### Post Booking

- Endpoints: `/booking`
- Method: `POST`
- Description: Post booking for a car
- Request Body:
```json
{
    "Date": "2024-06-01",
    "StartTime": "09:00:00",
    "EndTime": "18:00:00",
    "UserID": 1,
    "CarID": 1,
    "Model": "Toyota Camry",
    "PaymentID": 1
}
```
- Response:
    - `200 OK`: Returns successful message.
    - `500 Internal Server Error`: Failed to post booking

---

### Check Booking Validity

- Endpoints: `/checkValidity`
- Method: `PUT`
- Description: Checks if a booking is valid
- Request Body:
```json
{
    "Date": "2024-07-01",
    "StartTime": "10:00:00", 
    "EndTime": "11:00:00",
    "CarID": 3
}
```
- Response:
    - `200 OK`: Returns successful message.
    ```json
    {
        "StatusCode": 200,
        "ResultMessage": "Booking is Valid."
    }
    ```
    - `401 Unauthorised`: Booking is not valid
    ```json
    {
        "StatusCode": 401,
        "ResultMessage": "Booking is Not Valid."
    }
    ```
---

### Get Booking 

- Endpoints: `/booking/{id}`
- Method: `GET`
- Description: Gets booking's info 
- Path Parameter: 
    - `id`: The unique id for each booking.
- Response:
    - `200 OK`: Returns booking information in JSON format.
    ```json
        {
            "BookingID": 1,
            "StartTime": "09:00:00",
            "EndTime": "17:00:00",
            "Date": "2024-12-20",
            "CarID": 1,
            "Model": "Toyota Camry",
            "UserID": 1,
            "PaymentID": 1
        }
    ```
    - `500 Internal Server Error`: Failed to retrieve booking information

---

### Update Booking 

- Endpoints: `/booking/{id}`
- Method: `PUT`
- Description: Update information of booking
- Path Parameter: 
    - `id`: The unique identifier of each booking.
- Request Body:
```json
{
    "StartTime": "021:00:00",
    "EndTime": "17:00:00",
    "Date": "2024-06-01"
}
```
- Response:
    - `200 OK`: Booking has been updated.
    - `500 Internal Server Error`: An error has occured while updating booking

---

### Delete Booking 

- Endpoints: `/booking/{id}`
- Method: `PUT`
- Description: Deletes booking
- Path Parameter: 
    - `id`: The unique identifier of each booking.
- Response:
    - `200 OK`: Booking has been deleted.
    - `500 Internal Server Error`: An error has occured while deleting booking

---

### Get Booking By Car ID

- Endpoints: `/booking/car/{id}`
- Method: `GET`
- Description: Gets all car's booking by car id
- Path Parameter: 
    - `id`: The unique id for each car.
- Response:
    - `200 OK`: Returns all bookings information in JSON format.
    ```json
    [
        {
            "BookingID": 2,
            "StartTime": "10:00:00",
            "EndTime": "16:00:00",
            "Date": "2024-12-20",
            "CarID": 2,
            "Model": "Honda Civic",
            "UserID": 2,
            "PaymentID": 2
        }, (...and more)
    ]
    ```
    - `500 Internal Server Error`: Failed to retrieve cars' booking information

---

### Get Booking By User ID

- Endpoints: `/booking/user/{id}`
- Method: `GET`
- Description: Gets all user's bookings by user id
- Path Parameter: 
    - `id`: The unique id for each user.
- Response:
    - `200 OK`: Returns all bookings information in JSON format.
    ```json
    [
        {
            "BookingID": 2,
            "StartTime": "10:00:00",
            "EndTime": "16:00:00",
            "Date": "2024-12-20",
            "CarID": 2,
            "Model": "Honda Civic",
            "UserID": 2,
            "PaymentID": 2
        }, (...and more)
    ]
    ```
    - `500 Internal Server Error`: Failed to retrieve user's booking information

---

## Return

[Click on this link to return to parent README.md](../../README.md)