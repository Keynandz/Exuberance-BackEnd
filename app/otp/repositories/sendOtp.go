package otp

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

const (
	CONFIG_SMTP_HOST   = "smtp.gmail.com"
	CONFIG_SMTP_PORT   = 587
	CONFIG_SENDER_NAME = "EXUBERANCE | Indonesia"
)

func Otp(recipientEmail, otpCode string) error {
	loadErr := godotenv.Load()
	if loadErr != nil {
		log.Fatal("error loading file .env")
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME+" <"+os.Getenv("GOOGLE_EMAIL")+">")
	mailer.SetHeader("To", recipientEmail)
	mailer.SetHeader("Subject", "OTP CODE REGISTRATION")
	bodyHTML := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>OTP Code Registration</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f0f0f0;
				}
				.container {
					margin: 20px auto;
					max-width: 500px;
					padding: 20px;
					background-color: #ffffff;
					border-radius: 8px;
					box-shadow: 0 0 8px rgba(0, 0, 0, 0.2);
				}
				.header {
					text-align: center;
					margin-bottom: 20px;
				}
				.header h1 {
					color: #333333;
				}
				.otp-code {
					text-align: center;
					font-size: 28px;
					font-weight: bold;
					color: #007bff;
					padding: 20px;
					background-color: #f0f0f0;
					border-radius: 8px;
				}
				.footer {
					text-align: center;
					margin-top: 20px;
					color: #777777;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>OTP Code Registration</h1>
				</div>
				<div class="otp-code">Your OTP: %s</div>
				<div class="footer">
					<p>This email is auto-generated. Please do not reply.</p>
				</div>
			</div>
		</body>
		</html>
	`, otpCode)
	mailer.SetBody("text/html", bodyHTML)

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		os.Getenv("GOOGLE_EMAIL"),
		os.Getenv("GOOGLE_PASS"),
	)

	if err := dialer.DialAndSend(mailer); err != nil {
		return err
	}

	log.Println("OTP email sent!")
	return nil
}
