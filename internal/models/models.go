package models

import "time"

// User is the user model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Room is the rooms models
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Room is the rooms models
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
}

type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	RoomID        int
	Reservation   Reservation
	Restriction   Restriction
}
