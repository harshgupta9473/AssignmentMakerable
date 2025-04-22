package utils

import (
	"net/smtp"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func SendVerificationEmail(email, token string, userId int64) error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	from := os.Getenv("emailID")
	password := os.Getenv("apppassword")
	smtpHost := os.Getenv("smtpHost")
	smtpPort := os.Getenv("smtpPort")
	userIdStr := strconv.FormatInt(userId, 10)
	link := os.Getenv("link") + "/verify?token=" + token + "&userID=" + userIdStr

	msg := "From: " + from + "\n" +
		"To: " + email + "\n" +
		"Subject: Email Verification \n\n" +
		"Click the link to verify your email: " + link

	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, []byte(msg))

}
