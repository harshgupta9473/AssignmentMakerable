package models

import "time"

type Patient struct {
	ID              int64     `json:"id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	DOB             time.Time `json:"dob"`
	Gender          string    `json:"gender"`
	ComplaintType   string    `json:"complaint_type"`
	DoctorDiagnosis string    `json:"doctor_diagnosis,omitempty"`
	DoctorID        int64     `json:"doctor_id"` // Added field
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type PatientRequest struct {
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	DOB             time.Time `json:"dob"`
	Gender          string    `json:"gender"`
	ComplaintType   string    `json:"complaint_type"`
	DoctorDiagnosis string    `json:"doctor_diagnosis,omitempty"`
	DoctorId        int64     `json:"doctor_id"`
}
