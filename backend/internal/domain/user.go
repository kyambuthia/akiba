package domain

import "time"

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusDisabled UserStatus = "disabled"
)

type User struct {
	ID            string
	EmailLower    string
	PhoneE164     string
	UsernameLower string
	PasswordHash  string
	Status        UserStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
