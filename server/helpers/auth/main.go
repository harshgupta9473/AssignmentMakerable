package authHelper

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/harshgupta9473/assignment_makerable/db"
	"github.com/harshgupta9473/assignment_makerable/models"
)

// ISUserExistsByEmail checks if a user exists by email and role.
func ISUserExistsByEmailandRole(email string, role string) (bool, error) {
	conn := db.GetDB()
	query := `SELECT 1 FROM users WHERE email=$1 AND role=$2`

	var exists int
	err := conn.QueryRow(query, email, role).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("error checking if user exists: %v", err)
	}
	return true, nil
}

// ISUserExistsByEmail checks if a user exists by email.
func ISUserExistsByEmail(email string) (bool, error) {
	conn := db.GetDB()
	query := `SELECT 1 FROM users WHERE email=$1`

	var exists int
	err := conn.QueryRow(query, email).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("error checking if user exists: %v", err)
	}
	return true, nil
}

// InsertDoctorWithGovID inserts a doctor into both `users` and `doctorgovid` using a transaction.
func InsertDoctorWithGovID(doctor models.SignupRequestDoctor, verificationToken string, expiryTime time.Time) (int64, error) {
	conn := db.GetDB()

	tx, err := conn.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %v", err)
	}

	var userID int64

	insertUserQuery := `
		INSERT INTO users (first_name, last_name, email, password, role, verification_token, token_expiry)
		VALUES ($1, $2, $3, $4, 'doctor', $5, $6)
		RETURNING id
	`

	err = tx.QueryRow(insertUserQuery,
		doctor.FirstName, doctor.LastName, doctor.Email, doctor.Password, verificationToken, expiryTime).
		Scan(&userID)

	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert doctor into users table: %v", err)
	}

	insertGovIDQuery := `
		INSERT INTO doctorgovid (gov_id, user_id)
		VALUES ($1, $2)
	`
	_, err = tx.Exec(insertGovIDQuery, doctor.DoctorID, userID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert into doctorgovid: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return userID, nil
}

func InsertIntoReceptionist(user models.SignupRequestReceptionist, verificationToken string, expiryTime time.Time) (int64, error) {
	query := `INSERT INTO users (first_name, last_name, email, password, role, verification_token, token_expiry)
			  VALUES ($1, $2, $3, $4, 'receptionist', $5, $6) RETURNING id`

	var userID int64
	err := db.GetDB().QueryRow(query,
		user.FirstName, user.LastName, user.Email, user.Password, verificationToken, expiryTime).
		Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert receptionist: %v", err)
	}

	return userID, nil
}

func UpdateUserVerificationToken(userID int64, token string, expiry time.Time) error {
	query := `UPDATE users SET verification_token = $1, token_expiry = $2 WHERE id = $3`
	_, err := db.GetDB().Exec(query, token, expiry, userID)
	return err
}

func GetUserByUserId(userID int64) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, email, role, is_verified, is_approved, verified_at, verification_token, token_expiry
		FROM users
		WHERE id = $1
	`

	row := db.GetDB().QueryRow(query, userID)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Role,
		&user.IsVerified,
		&user.IsApproved,
		&user.VerifiedAt,
		&user.VerificationToken,
		&user.TokenExpiry,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}


func GetUserByUserByEmailAndRole(email string, role string) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, email,password, role, is_verified, is_approved, verified_at, verification_token, token_expiry
		FROM users
		WHERE email = $1 AND role = $2
	`

	row := db.GetDB().QueryRow(query, email, role)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsVerified,
		&user.IsApproved,
		&user.VerifiedAt,
		&user.VerificationToken,
		&user.TokenExpiry,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func IsUserApproved(userID int64) (bool, error) {
	var isApproved bool
	err := db.GetDB().QueryRow(
		`SELECT is_approved FROM users WHERE id = $1`, userID,
	).Scan(&isApproved)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("doctor not found")
		}
		return false, err
	}
	return isApproved, nil
}


func GetTokenExpiryIsVerified(userID int64)(string,time.Time,bool,error){
	var storedToken string
	var expiry time.Time
	var isVerified bool

	query := `SELECT verification_token, token_expiry, is_verified FROM users WHERE id = $1`
	err := db.GetDB().QueryRow(query, userID).Scan(&storedToken, &expiry, &isVerified)
	if err!=nil{
		return "",time.Time{},false,err
	}
	return storedToken,expiry,isVerified,nil

}


func VerifyUserwithUserID(userID int64)error{
	updateQuery := `UPDATE users SET is_verified = TRUE, verified_at = $1 WHERE id = $2`
	_, err := db.GetDB().Exec(updateQuery, time.Now(), userID)
	return err
}