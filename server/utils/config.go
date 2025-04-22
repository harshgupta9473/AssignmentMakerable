package utils

import "os"

var (
	JwtSecretKey string
	DocotorRole string="doctor"
	ReceptionistRole string="receptionist"
	AdminRole string="admin"
)

func LoadSecrets() {
	JwtSecretKey = os.Getenv("JWTSecretKey");
}