package main

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var port int = 8000

type User struct {
	UserID         int
	Name           string
	EmailAddr      string
	ContactNo      string
	MembershipTier string
	DateJoined     time.Time
	asswordHash    string
}
