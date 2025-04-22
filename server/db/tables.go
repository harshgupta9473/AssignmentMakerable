package db

import (
	"log"
)

// CreateAllTable creates all  tables for the application
func CreateAllTable() error {
	err:=CreateAdminTable()
	if err!=nil{
		log.Println("Error creating admin table")
		return err
	}

	err = CreateUsersTable()
	if err != nil {
		log.Println("Error creating users table:", err)
		return err
	}

	err = CreateTableForDoctorsGovernmentID()
	if err != nil {
		log.Println("Error creating doctorgovid table:", err)
		return err
	}

	err = CreatePatientsTable()
	if err != nil {
		log.Println("Error creating patients table:", err)
		return err
	}

	return nil
}
func CreateAdminTable()error{
	query:=`CREATE TABLE IF NOT EXISTS admins (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`
   _,err:=dB.Exec(query)
   return err
}
func CreateUsersTable() error {

	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			role TEXT CHECK(role IN ('doctor', 'receptionist','admin')) NOT NULL,
			is_verified BOOLEAN DEFAULT FALSE,
			is_approved BOOLEAN DEFAULT FALSE,
			verified_at TIMESTAMP,
			verification_token TEXT,
			token_expiry TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := dB.Exec(query)
	if err != nil {
		log.Println("Error creating users table:", err)
		return err
	}
	return nil
}

func CreateTableForDoctorsGovernmentID() error {

	query := `
		CREATE TABLE IF NOT EXISTS doctorgovid (
			id SERIAL PRIMARY KEY,
			gov_id VARCHAR(50) NOT NULL,
			user_id INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`

	_, err := dB.Exec(query)
	if err != nil {
		log.Println("Error creating doctorgovid table:", err)
		return err
	}
	return nil
}

func CreatePatientsTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS patients (
        id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL,
    dob DATE NOT NULL,
    gender VARCHAR(20),
    complaint_type VARCHAR(255),
    doctor_diagnosis TEXT,
    doctor_id INTEGER,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (doctor_id) REFERENCES users(id) ON DELETE SET NULL
    );
    `
	_, err := dB.Exec(query)
	return err
}
