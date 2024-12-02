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
	IsActivated TINYINT(1) NOT NULL, 
    VerificationCodeHash VARCHAR(100) NOT NULL
);

-- Car Table 
CREATE TABLE Car (
	CarID INT PRIMARY KEY AUTO_INCREMENT,
    Model VARCHAR(20) NOT NULL,
    PlateNo VARCHAR(10) NOT NULL, 
    RentalRate INT NOT NULL
);

-- Payment Table
CREATE TABLE Payment (
	PaymentID INT PRIMARY KEY AUTO_INCREMENT,
    Amount DECIMAL(10, 2) NOT NULL, 
    DateCreated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    Status ENUM('Pending', 'Successful', 'Refunded', 'Unsuccessful') NOT NULL DEFAULT 'Pending',
    UserID INT NOT NULL, 
    CarID INT NOT NULL,
    
    FOREIGN KEY(UserID) REFERENCES User(UserID),
    FOREIGN KEY(CarID) REFERENCES Car(CarID)
);

-- Booking Table
CREATE TABLE Booking (
	BookingID INT PRIMARY KEY AUTO_INCREMENT, 
    StartDate DATE NOT NULL,
    EndDate DATE NOT NULL,
    StartTime TIME NOT NULL,
    EndTime TIME NOT NULL,
    UserID INT NOT NULL,
    CarID INT NOT NULL, 
    PaymentID INT NOT NULL, 
    
	FOREIGN KEY(UserID) REFERENCES User(UserID),
    FOREIGN KEY(CarID) REFERENCES Car(CarID),
    FOREIGN KEY(PaymentID) REFERENCES Payment(PaymentID)
);

-- Data Creation

-- Inserting data into User table
INSERT INTO User (Name, EmailAddr, ContactNo, MemberTier, PasswordHash, IsActivated, VerificationCodeHash)
VALUES 
('Caden Toh', 'cadentohjunyi@gmail.com', '84469588', 'Basic', '$2a$10$WMkzkeV/CroPDMPwrk8Q4ONTN7wh71K0ObS.KypcCF541lwaRwm3a', 1, 'hash1'),
('John Doe', 'john@example.com', '12345678', 'Premium', 'hashed_password1', 1, 'hash2'),
('Jane Smith', 'jane@example.com', '87654321', 'Basic', 'hashed_password2', 1, 'hash3'),
('Michael Jones', 'michael@example.com', '98765432', 'VIP', 'hashed_password3', 1, 'hash4');

-- Inserting data into Car table
INSERT INTO Car (Model, PlateNo, RentalRate)
VALUES 
('Toyota Camry', 'ABC1234', 50),
('Honda Civic', 'DEF5678', 40),
('Tesla Model 3', 'GHI987', 70);

-- Inserting data into Payment table
INSERT INTO Payment (Amount, Status, UserID, CarID)
VALUES 
(200, 'Pending', 1, 1),
(150, 'Successful', 2, 2),
(250, 'Successful', 3, 3);

-- Inserting data into Booking table
INSERT INTO Booking (StartDate, EndDate, StartTime, EndTime, UserID, CarID, PaymentID)
VALUES 
('2024-06-01', '2024-06-03', '09:00:00', '17:00:00', 1, 1, 1),
('2024-06-15', '2024-06-16', '10:00:00', '16:00:00', 2, 2, 2),
('2024-07-01', '2024-07-02', '11:00:00', '15:00:00', 3, 3, 3);