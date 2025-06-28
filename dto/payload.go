package dto

import "time"

type CreateLoanPayload struct {
	BorrowerID      string  `json:"borrower_id" binding:"required"`
	PrincipalAmount float64 `json:"principal_amount" binding:"required,gt=0"`
	Rate            float64 `json:"rate" binding:"required,gt=0"`
	ROI             float64 `json:"roi" binding:"required,gte=0"`
}

type ApproveLoanPayload struct {
	LoanID       int       `json:"-"`
	PictureProof string    `json:"picture_proof" binding:"required"`
	EmployeeID   string    `json:"employee_id" binding:"required"`
	DateStr      string    `json:"date" binding:"required"`
	Date         time.Time `json:"-"`
}

type InvestLoanPayload struct {
	LoanID        int     `json:"-"`
	InvestorEmail string  `json:"investor_email" binding:"required,email"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
}

type DisburseLoanPayload struct {
	LoanID        int       `json:"-"`
	AgreementLink string    `json:"agreement_letter_link" binding:"required,url"`
	EmployeeID    string    `json:"employee_id" binding:"required"`
	DateStr       string    `json:"date" binding:"required"`
	Date          time.Time `json:"-"`
}
