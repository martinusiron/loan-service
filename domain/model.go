package domain

import "time"

type LoanStatus string

const (
	StatusProposed  LoanStatus = "proposed"
	StatusApproved  LoanStatus = "approved"
	StatusInvested  LoanStatus = "invested"
	StatusDisbursed LoanStatus = "disbursed"
)

type Loan struct {
	ID                  int
	BorrowerID          string
	PrincipalAmount     float64
	Rate                float64
	ROI                 float64
	Status              LoanStatus
	AgreementLetterLink string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type LoanApproval struct {
	ID           int
	LoanID       int
	PictureProof string
	EmployeeID   string
	ApprovedAt   time.Time
}

type Investment struct {
	ID            int
	LoanID        int
	InvestorEmail string
	Amount        float64
	InvestedAt    time.Time
}
