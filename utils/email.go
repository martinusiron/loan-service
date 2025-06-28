package utils

import (
	"fmt"
	"log"
)

func SendDummyEmail(to string, loanID int) {
	subject := fmt.Sprintf("Loan #%d has been fully funded", loanID)
	body := "Thank you for your investment. The agreement letter will be sent soon."

	// Simulate email sending with a log
	log.Printf("[EMAIL SENT] To: %s | Subject: %s | Body: %s", to, subject, body)
}
