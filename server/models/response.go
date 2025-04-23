package models

import "time"

type DoctorResponse struct {
	ID         int64      `json:"id"`
	GovID      string     `json:"gov_id"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Role       string      `json:"role"`
	Email      string     `json:"email,omitempty"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
}
