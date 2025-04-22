package models

import "time"

type SignupRequestDoctor struct {
	DoctorID  string `json:"doctor_id"` // Allotted by government government ID alloted to doctors
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type SignupRequestReceptionist struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type SignUpRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type User struct {
	ID                int64      `json:"id"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	Email             string     `json:"email"`
	Role              string     `json:"role"`
	Password          string     `json:"-"`
	IsVerified        bool       `json:"is_verified"`
	IsApproved        bool       `json:"is_approved,omitempty"`
	VerifiedAt        *time.Time `json:"verified_at,omitempty"`
	VerificationToken string     `json:"verification_token"`
	TokenExpiry       time.Time  `json:"token_expiry"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
