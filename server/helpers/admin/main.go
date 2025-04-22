package adminHelper

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/harshgupta9473/assignment_makerable/db"
	"github.com/harshgupta9473/assignment_makerable/models"
)

func EmailExistsInAdmins(email string) (bool, error) {
	conn := db.GetDB()

	var exists bool
	selectQuery := `SELECT EXISTS(SELECT 1 from admins WHERE email = $1)`
	err := conn.QueryRow(selectQuery, email).Scan(&exists)
	if err != nil {
		log.Println("Error checking email in admins table:", err)
		return false, err
	}

	return exists, nil
}

func CreateAdminWithVerification(user models.SignUpRequest, verificationToken string, expiryTime time.Time) (int64, error) {
	query := `INSERT INTO users (first_name, last_name, email, password, role, verification_token, token_expiry)
			  VALUES ($1, $2, $3, $4, 'admin', $5, $6) RETURNING id`

	var userID int64
	err := db.GetDB().QueryRow(query,
		user.FirstName, user.LastName, user.Email, user.Password, verificationToken, expiryTime).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert receptionist: %v", err)
	}
	return userID, nil
}

func GetUsersByApproval(approved string,role string) ([]models.User, error) {
	db := db.GetDB()
	query := `SELECT id, first_name, last_name, email, role, is_verified, is_approved, verified_at, verification_token, token_expiry
		FROM users
		WHERE role=$1`
	var rows *sql.Rows
	var err error

	if approved == "true" {
		query += `AND is_approved = true`
	} else if approved == "false" {
		query += `AND is_approved = false`
	}

	rows, err = db.Query(query,role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Role,
			&user.IsVerified,
			&user.IsApproved,
			&user.VerifiedAt,
			&user.VerificationToken,
			&user.TokenExpiry,); err != nil {
			return nil, err
		}
		users = append(users,user)
	}
	return users, nil
}


func ApproveUserByRole(userId int64, role string) error {
	query := `
		UPDATE users
		SET is_approved = true
		WHERE id = $1 AND role = $2
	`
	_, err := db.GetDB().Exec(query, userId, role)
	return err
}
