DROP DATABASE IF EXISTS CNADAssg1DB;
CREATE DATABASE CNADAssg1DB;

USE CNADAssg1DB;

-- Table Creation 

-- User Table 
CREATE TABLE User (
	UserID INT PRIMARY KEY AUTO_INCREMENT,
    Name VARCHAR(50) NOT NULL,
    EmailAddr VARCHAR(100) UNIQUE NOT NULL, 
    ContactNo VARCHAR(8) UNIQUE NOT NULL CHECK (ContactNo NOT LIKE '%[^0-9]%'),
    MemberTier ENUM('Basic', 'Premium', 'VIP') NOT NULL DEFAULT 'Basic',
	PasswordHash VARCHAR(100) NOT NULL,
	IsActivated TINYINT(1) NOT NULL DEFAULT 0, 
    VerificationCodeHash VARCHAR(100) DEFAULT 'NoHash'
);

-- Car Table 
CREATE TABLE Car (
	CarID INT PRIMARY KEY AUTO_INCREMENT,
    Model VARCHAR(20) NOT NULL,
    PlateNo VARCHAR(10) NOT NULL, 
    RentalRate INT NOT NULL,
    Location VARCHAR(100) NOT NULL
);

CREATE TABLE Promotion (
	PromotionCode INT PRIMARY KEY AUTO_INCREMENT,
    Name VARCHAR(50) NOT NULL,
    Description VARCHAR(100) NOT NULL, 
    Discount DECIMAL(10,2) NOT NULL
);

-- Payment Table
CREATE TABLE Payment (
	PaymentID INT PRIMARY KEY AUTO_INCREMENT,
    Amount DECIMAL(10, 2) NOT NULL, 
    DateCreated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    Status ENUM('Pending', 'Successful', 'Refunded', 'Unsuccessful') NOT NULL DEFAULT 'Pending',
    PromotionCode INT NULL,
    UserID INT NOT NULL, 
    CarID INT NOT NULL,
    
    FOREIGN KEY(PromotionCode) REFERENCES Promotion(PromotionCode),
    FOREIGN KEY(UserID) REFERENCES User(UserID),
    FOREIGN KEY(CarID) REFERENCES Car(CarID)
);

-- Booking Table
CREATE TABLE Booking (
	BookingID INT PRIMARY KEY AUTO_INCREMENT, 
    Date DATE NOT NULL,
    StartTime TIME NOT NULL,
    EndTime TIME NOT NULL,
    UserID INT NOT NULL,
    CarID INT NOT NULL, 
    Model VARCHAR(20) NOT NULL,
    PaymentID INT NOT NULL, 
    
	FOREIGN KEY(UserID) REFERENCES User(UserID),
    FOREIGN KEY(CarID) REFERENCES Car(CarID),
    FOREIGN KEY(PaymentID) REFERENCES Payment(PaymentID)
);

-- Data Creation
-- Inserting data into User table
INSERT INTO User (Name, EmailAddr, ContactNo, MemberTier, PasswordHash, IsActivated, VerificationCodeHash)
VALUES 
('John Doe', 'johndoe@gmail.com', '21212121', 'Premium', 'hashed_password1', 1, 'hash2'),
('Cassie', 'unknownsentinel08@gmail.com', '12345678', 'Premium', '$2a$10$hoEOy0Me9RNwxsFyTx2/S.WsNPV0WmS9jCAAwe6zDpiAs9NfwmwZK', 1, 'hash2'),
('Jane Smith', 'jane@example.com', '87654321', 'Basic', 'hashed_password2', 1, 'hash3'),
('Michael Jones', 'michael@example.com', '98765432', 'VIP', 'hashed_password3', 1, 'hash4');

-- Inserting data into Car table
INSERT INTO Car (Model, PlateNo, RentalRate, Location)
VALUES 
('Toyota Camry', 'ABC1234', 50, 'Lorong Ah Soo'),
('Honda Civic', 'DEF5678', 40, 'Lorong Ah Soo'),
('BMW Type X', 'DEF5678', 40, 'Lorong Ah Soo'),
('Tesla Model 3', 'GHI987', 70, 'Lorong Ah Soo');

-- Inserting data into Promotion table
INSERT INTO Promotion (Name, Description, Discount)
VALUES 
('Christmas Discount', 'A limited time offer discount in lieu of the upcoming festivities', 0.1),
('Premium Member Discount', 'DEF5678', 0.05),
('VIP Member Discount', 'GHI987', 0.1);

-- Inserting data into Payment table
INSERT INTO Payment (Amount, Status, PromotionCode, UserID, CarID)
VALUES 
(200, 'Pending', NULL, 1, 1),
(150, 'Successful', NULL, 2, 2),
(250, 'Successful',NULL, 3, 3);

-- Inserting data into Booking table
INSERT INTO Booking (Date, StartTime, EndTime, UserID, CarID, Model, PaymentID)
VALUES 
('2024-06-01', '09:00:00', '17:00:00', 1, 1, 'Toyota Camry', 1),
('2024-06-15', '10:00:00', '16:00:00', 2, 2, 'Honda Civic', 2),
('2024-07-01', '11:00:00', '15:00:00', 3, 3, 'BMW Type X', 3);

DELIMITER //
CREATE TRIGGER updatePaymentOnBookingDelete
AFTER DELETE ON Booking
FOR EACH ROW
BEGIN
    UPDATE Payment
    SET Status = 'Refunded'
    WHERE PaymentID = OLD.PaymentID;
END;
//
DELIMITER ;

DELIMITER //

CREATE PROCEDURE CheckBookingValidity (
    IN newDate DATE,
    IN newStartTime TIME,
    IN newEndTime TIME,
    IN newCarID INT,
    OUT statusCode INT,
    OUT resultMessage VARCHAR(255)
)
BEGIN
    -- Declare a variable to count conflicting bookings
    DECLARE conflictCount INT;

    -- Check for time conflicts on the same date and car
    SELECT COUNT(*)
    INTO conflictCount
    FROM Booking
    WHERE 
        Date = newDate 
        AND CarID = newCarID
		AND (
            -- Overlapping time conditions, excluding exact end-to-start or start-to-end cases
            (newStartTime < EndTime AND newStartTime >= StartTime) 
            OR (newEndTime > StartTime AND newEndTime <= EndTime)
            OR (newStartTime < StartTime AND newEndTime > EndTime)
        );

    -- If there's a conflict, set the result message and do not insert the booking
    IF conflictCount > 0 THEN
		SET statusCode = 401;
        SET resultMessage = 'Booking failed: Time slot conflicts with an existing booking.';
    ELSE
		SET statusCode = 200;
        SET resultMessage = 'Booking is Valid.';
    END IF;
    
    SELECT statusCode, resultMessage;
END;
//

DELIMITER ;
