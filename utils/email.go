package utils

import "log"

func SendEmail(to, subject, body string) {
	log.Printf("Sending email to %s | Subject: %s | Body: %s\n", to, subject, body)
}
