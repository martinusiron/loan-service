package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/martinusiron/loan-service/domain"

	mockRepo "github.com/martinusiron/loan-service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateLoan(t *testing.T) {
	mockLoanRepo := new(mockRepo.LoanRepository)
	mockApprovalRepo := new(mockRepo.ApprovalRepository)
	mockInvestRepo := new(mockRepo.InvestmentRepository)

	uc := NewLoanUsecase(mockLoanRepo, mockApprovalRepo, mockInvestRepo)

	mockLoanRepo.On("CreateLoan", mock.Anything, mock.AnythingOfType("*domain.Loan")).Return(nil)

	loan, err := uc.CreateLoan(context.TODO(), "BR123", 1000000, 10.0, 5.0)

	assert.NoError(t, err)
	assert.Equal(t, "BR123", loan.BorrowerID)
	assert.Equal(t, float64(1000000), loan.PrincipalAmount)
}

func TestApproveLoan(t *testing.T) {
	mockLoanRepo := new(mockRepo.LoanRepository)
	mockApprovalRepo := new(mockRepo.ApprovalRepository)
	mockInvestRepo := new(mockRepo.InvestmentRepository)

	uc := NewLoanUsecase(mockLoanRepo, mockApprovalRepo, mockInvestRepo)

	loan := &domain.Loan{
		ID:         1,
		Status:     domain.StatusProposed,
		BorrowerID: "B01",
	}
	mockLoanRepo.On("GetLoanByID", mock.Anything, 1).Return(loan, nil)
	mockApprovalRepo.On("CreateApproval", mock.Anything, mock.AnythingOfType("*domain.LoanApproval")).Return(nil)
	mockLoanRepo.On("UpdateLoanStatus", mock.Anything, 1, domain.StatusApproved).Return(nil)

	err := uc.ApproveLoan(context.TODO(), 1, "proof.jpg", "EMP001", time.Now())
	assert.NoError(t, err)
}

func TestInvestLoan_Full(t *testing.T) {
	mockLoanRepo := new(mockRepo.LoanRepository)
	mockApprovalRepo := new(mockRepo.ApprovalRepository)
	mockInvestRepo := new(mockRepo.InvestmentRepository)

	uc := NewLoanUsecase(mockLoanRepo, mockApprovalRepo, mockInvestRepo)

	loan := &domain.Loan{
		ID:              1,
		Status:          domain.StatusApproved,
		PrincipalAmount: 1000000,
	}

	mockLoanRepo.On("GetLoanByID", mock.Anything, 1).Return(loan, nil)
	mockInvestRepo.On("GetTotalInvested", mock.Anything, 1).Return(900000.0, nil)
	mockInvestRepo.On("AddInvestment", mock.Anything, mock.AnythingOfType("*domain.Investment")).Return(nil)
	mockLoanRepo.On("UpdateLoanStatus", mock.Anything, 1, domain.StatusInvested).Return(nil)
	mockInvestRepo.On("GetInvestorsByLoan", mock.Anything, 1).Return([]domain.Investment{
		{InvestorEmail: "a@a.com", Amount: 500000},
		{InvestorEmail: "b@b.com", Amount: 500000},
	}, nil)

	err := uc.InvestLoan(context.TODO(), 1, "new@investor.com", 100000)
	assert.NoError(t, err)
}

func TestDisburseLoan(t *testing.T) {
	mockLoanRepo := new(mockRepo.LoanRepository)
	mockApprovalRepo := new(mockRepo.ApprovalRepository)
	mockInvestRepo := new(mockRepo.InvestmentRepository)

	uc := NewLoanUsecase(mockLoanRepo, mockApprovalRepo, mockInvestRepo)

	loan := &domain.Loan{
		ID:     1,
		Status: domain.StatusInvested,
	}
	mockLoanRepo.On("GetLoanByID", mock.Anything, 1).Return(loan, nil)
	mockLoanRepo.On("SetAgreementLink", mock.Anything, 1, "http://link.com/file.pdf").Return(nil)
	mockLoanRepo.On("UpdateLoanStatus", mock.Anything, 1, domain.StatusDisbursed).Return(nil)

	err := uc.DisburseLoan(context.TODO(), 1, "http://link.com/file.pdf", "EMP02", time.Now())
	assert.NoError(t, err)
}
