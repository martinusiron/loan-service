package repository

import (
	"context"

	"github.com/martinusiron/loan-service/domain"
)

type LoanRepository interface {
	CreateLoan(ctx context.Context, loan *domain.Loan) error
	GetLoanByID(ctx context.Context, id int) (*domain.Loan, error)
	UpdateLoanStatus(ctx context.Context, id int, status domain.LoanStatus) error
	SetAgreementLink(ctx context.Context, id int, link string) error
}

type ApprovalRepository interface {
	CreateApproval(ctx context.Context, a *domain.LoanApproval) error
}

type InvestmentRepository interface {
	AddInvestment(ctx context.Context, i *domain.Investment) error
	GetTotalInvested(ctx context.Context, loanID int) (float64, error)
	GetInvestorsByLoan(ctx context.Context, loanID int) ([]domain.Investment, error)
}
