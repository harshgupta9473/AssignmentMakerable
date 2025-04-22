package patientsHelper

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/harshgupta9473/assignment_makerable/db"
	"github.com/harshgupta9473/assignment_makerable/models"
)

func CreatePatient(patient models.PatientRequest) (int64, error) {
	var id int64
	query := `INSERT INTO patients (first_name, last_name, email, phone, dob, gender, complaint_type, doctor_diagnosis, doctor_id, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`
	err := db.GetDB().QueryRow(query, patient.FirstName, patient.LastName, patient.Email, patient.Phone, patient.DOB, patient.Gender, patient.ComplaintType, patient.DoctorDiagnosis, patient.DoctorId, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("could not create patient: %v", err)
	}
	return id, nil
}

func UpdatePatient(patient models.Patient) error {
	if patient.ID == 0 {
		return fmt.Errorf("invalid patient ID")
	}
	query := `
		UPDATE patients
		SET first_name = $1, last_name = $2, email = $3, phone = $4, dob = $5, gender = $6, complaint_type = $7, doctor_diagnosis = $8, doctor_id = $9, updated_at = $10
		WHERE id = $11
	`

	_, err := db.GetDB().Exec(query, patient.FirstName, patient.LastName, patient.Email, patient.Phone, patient.DOB, patient.Gender, patient.ComplaintType, patient.DoctorDiagnosis, patient.DoctorID, time.Now(), patient.ID)
	if err != nil {
		return fmt.Errorf("could not update patient: %v", err)
	}

	return nil
}

func GetPatient(id int64) (models.Patient, error) {
	var patient models.Patient
	query := `SELECT id, first_name, last_name, email, phone, dob, gender, complaint_type, doctor_diagnosis, created_at, updated_at FROM patients WHERE id = $1 AND is_deleted=false`
	err := db.GetDB().QueryRow(query, id).Scan(&patient.ID, &patient.FirstName, &patient.LastName, &patient.Email, &patient.Phone, &patient.DOB, &patient.Gender, &patient.ComplaintType, &patient.DoctorDiagnosis, &patient.CreatedAt, &patient.UpdatedAt)
	if err != nil {
		return models.Patient{}, fmt.Errorf("could not retrieve patient: %v", err)
	}
	return patient, nil
}

func UpdateDoctorDiagnosis(patientID int, doctorDiagnosis string) error {
	if patientID == 0 {
		return fmt.Errorf("invalid patient ID")
	}

	query := `
		UPDATE patients
		SET doctor_diagnosis = $1, updated_at = $2
		WHERE id = $3 AND is_deleted=false
	`

	_, err := db.GetDB().Exec(query, doctorDiagnosis, time.Now(), patientID)
	if err != nil {
		return fmt.Errorf("could not update doctor diagnosis: %v", err)
	}

	return nil
}

func GetPatientsByDoctorID(doctorID int64) ([]models.Patient, error) {

	if doctorID == 0 {
		return nil, fmt.Errorf("invalid doctor ID")
	}

	query := `
		SELECT id, first_name, last_name, email, phone, dob, gender, complaint_type, doctor_diagnosis, created_at, updated_at
		FROM patients
		WHERE doctor_id = $1 AND is_deleted=false
	`

	var patients []models.Patient

	rows, err := db.GetDB().Query(query, doctorID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch patients for doctor ID %d: %v", doctorID, err)
	}
	defer rows.Close()

	for rows.Next() {
		var patient models.Patient
		if err := rows.Scan(&patient.ID, &patient.FirstName, &patient.LastName, &patient.Email, &patient.Phone, &patient.DOB, &patient.Gender, &patient.ComplaintType, &patient.DoctorDiagnosis, &patient.CreatedAt, &patient.UpdatedAt); err != nil {
			return nil, fmt.Errorf("could not scan patient row: %v", err)
		}
		patients = append(patients, patient)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while fetching rows: %v", err)
	}

	return patients, nil
}

func GetPatientByDoctorIDAndPatientID(doctorID, patientID int64) (*models.Patient, error) {
	if doctorID == 0 || patientID == 0 {
		return nil, fmt.Errorf("invalid doctor ID or patient ID")
	}

	query := `
		SELECT id, first_name, last_name, email, phone, dob, gender, complaint_type, doctor_diagnosis, created_at, updated_at
		FROM patients
		WHERE doctor_id = $1 AND id = $2 AND is_deleted=false
	`

	var patient models.Patient

	row := db.GetDB().QueryRow(query, doctorID, patientID)
	err := row.Scan(&patient.ID, &patient.FirstName, &patient.LastName, &patient.Email, &patient.Phone, &patient.DOB, &patient.Gender, &patient.ComplaintType, &patient.DoctorDiagnosis, &patient.CreatedAt, &patient.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {

			return nil, fmt.Errorf("no patient found for doctor ID %d and patient ID %d", doctorID, patientID)
		}
		return nil, fmt.Errorf("could not fetch patient details: %v", err)
	}

	return &patient, nil
}

func SoftDeletePatientByID(patientID int64) error {
	if patientID == 0 {
		return fmt.Errorf("invalid patient ID")
	}

	query := `
		UPDATE patients
		SET is_deleted = true, updated_at = $1
		WHERE id = $2
	`

	_, err := db.GetDB().Exec(query, time.Now(), patientID)
	if err != nil {
		return fmt.Errorf("could not soft delete patient: %v", err)
	}

	return nil
}

func GetAllPatients() ([]models.Patient, error) {
	var patients []models.Patient
	query := `SELECT id, first_name, last_name, email, phone, dob, gender, complaint_type, doctor_diagnosis, doctor_id, created_at, updated_at FROM patients WHERE is_deleted = false`
	rows, err := db.GetDB().Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching patients: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var patient models.Patient
		if err := rows.Scan(
			&patient.ID,
			&patient.FirstName,
			&patient.LastName,
			&patient.Email,
			&patient.Phone,
			&patient.DOB,
			&patient.Gender,
			&patient.ComplaintType,
			&patient.DoctorDiagnosis,
			&patient.DoctorID,
			&patient.CreatedAt,
			&patient.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("could not scan patient row: %v", err)
		}
		patients = append(patients, patient)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while fetching rows: %v", err)
	}

	return patients, nil
}
