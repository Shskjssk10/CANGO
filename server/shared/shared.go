package shared

// Type Structures

type User struct {
	UserID               int
	Name                 string
	EmailAddr            string
	ContactNo            string
	MembershipTier       string
	PasswordHash         string
	IsActivated          int
	VerificationCodeHash string
}

type Car struct {
	CarID      int
	Model      string
	PlateNo    string
	RentalRate int
}

type Booking struct {
	BookingID int
	StartTime string
	EndTime   string
	StartDate string
	EndDate   string
	CarID     int
	UserID    int
	PaymentID int
}

type Payment struct {
	PaymentID   int
	Amount      int
	DateCreated string
	Status      string
	UserID      int
	CarID       int
}
