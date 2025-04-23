package doctorHelper

import (
	"github.com/harshgupta9473/assignment_makerable/db"
	"github.com/harshgupta9473/assignment_makerable/models"
)

func GetAllApprovedAndVerifiedDoctorsDetailByRole() ([]models.DoctorResponse, error) {
	query := `
		SELECT u.id, u.first_name, u.last_name, u.role, u.verified_at, dg.gov_id
		FROM users AS u
		JOIN doctorgovid AS dg ON dg.user_id = u.id
		WHERE u.role = 'doctor' AND u.is_verified = true AND u.is_approved = true
	`
	rows, err := db.GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() 

	var doctors []models.DoctorResponse
	for rows.Next() {
		var doc models.DoctorResponse
		err := rows.Scan(&doc.ID, &doc.FirstName, &doc.LastName, &doc.Role, &doc.VerifiedAt, &doc.GovID)
		if err != nil {
			return nil, err
		}
		doctors = append(doctors, doc)
	}

	// row iteration errors may be  // less knowledge
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return doctors, nil
}


func GetADoctorById(id int64)(*models.DoctorResponse,error){
	query := `
		SELECT u.id, u.first_name, u.last_name, u.role, u.verified_at, dg.gov_id
		FROM users AS u
		JOIN doctorgovid AS dg ON dg.user_id = u.id
		WHERE u.role = 'doctor' AND u.is_verified = true AND u.is_approved = true AND u.id=$1`
		var doc models.DoctorResponse
		err:=db.GetDB().QueryRow(query,id).Scan(&doc.ID, &doc.FirstName, &doc.LastName, &doc.Role, &doc.VerifiedAt, &doc.GovID)
		if err!=nil{
			return nil,err
		}
		return &doc,nil
}