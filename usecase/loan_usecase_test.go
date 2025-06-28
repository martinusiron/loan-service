package usecase

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/martinusiron/loan-service/domain"
	"github.com/martinusiron/loan-service/dto"

	mockRepo "github.com/martinusiron/loan-service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateLoan(t *testing.T) {
	mockLoanRepo := new(mockRepo.LoanRepository)
	mockApprovalRepo := new(mockRepo.ApprovalRepository)
	mockInvestRepo := new(mockRepo.InvestmentRepository)
	db := &sql.DB{}

	uc := NewLoanUsecase(mockLoanRepo, mockApprovalRepo, mockInvestRepo, db)

	mockLoanRepo.On("CreateLoan", mock.Anything, mock.AnythingOfType("*domain.Loan")).Return(nil)

	payload := dto.CreateLoanPayload{
		BorrowerID:      "BR123",
		PrincipalAmount: 1000000,
		Rate:            10.0,
		ROI:             5.0,
	}
	loan, err := uc.CreateLoan(context.TODO(), payload)

	assert.NoError(t, err)
	assert.Equal(t, "BR123", loan.BorrowerID)
	assert.Equal(t, float64(1000000), loan.PrincipalAmount)
}

func TestApproveLoan(t *testing.T) {
	mockLoanRepo := new(mockRepo.LoanRepository)
	mockApprovalRepo := new(mockRepo.ApprovalRepository)
	mockInvestRepo := new(mockRepo.InvestmentRepository)
	db := &sql.DB{}

	uc := NewLoanUsecase(mockLoanRepo, mockApprovalRepo, mockInvestRepo, db)

	loan := &domain.Loan{
		ID:         1,
		Status:     domain.StatusProposed,
		BorrowerID: "B01",
	}
	mockLoanRepo.On("GetLoanByID", mock.Anything, 1).Return(loan, nil)
	mockApprovalRepo.On("CreateApproval", mock.Anything, mock.AnythingOfType("*domain.LoanApproval")).Return(nil)
	mockLoanRepo.On("UpdateLoanStatus", mock.Anything, 1, domain.StatusApproved).Return(nil)

	payload := dto.ApproveLoanPayload{
		LoanID:       1,
		PictureProof: "proof.jpg",
		EmployeeID:   "EMP001",
		Date:         time.Now(),
	}
	err := uc.ApproveLoan(context.TODO(), payload)
	assert.NoError(t, err)
}

func TestInvestLoan_Full(t *testing.T) {
	mockLoanRepo := new(mockRepo.LoanRepository)
	mockApprovalRepo := new(mockRepo.ApprovalRepository)
	mockInvestRepo := new(mockRepo.InvestmentRepository)
	db := &sql.DB{}

	uc := NewLoanUsecase(mockLoanRepo, mockApprovalRepo, mockInvestRepo, db)

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

	payload := dto.InvestLoanPayload{
		LoanID:        1,
		InvestorEmail: "new@investor.com",
		Amount:        100000,
	}
	err := uc.InvestLoan(context.TODO(), payload)
	assert.NoError(t, err)
}

func TestDisburseLoan(t *testing.T) {
	mockLoanRepo := new(mockRepo.LoanRepository)
	mockApprovalRepo := new(mockRepo.ApprovalRepository)
	mockInvestRepo := new(mockRepo.InvestmentRepository)
	db := &sql.DB{}

	uc := NewLoanUsecase(mockLoanRepo, mockApprovalRepo, mockInvestRepo, db)

	loan := &domain.Loan{
		ID:     1,
		Status: domain.StatusInvested,
	}
	mockLoanRepo.On("GetLoanByID", mock.Anything, 1).Return(loan, nil)
	mockLoanRepo.On("SetAgreementLink", mock.Anything, 1, "http://link.com/file.pdf").Return(nil)
	mockLoanRepo.On("UpdateLoanStatus", mock.Anything, 1, domain.StatusDisbursed).Return(nil)

	payload := dto.DisburseLoanPayload{
		LoanID:        1,
		AgreementLink: "http://link.com/file.pdf",
		EmployeeID:    "EMP02",
		Date:          time.Now(),
	}
	err := uc.DisburseLoan(context.TODO(), payload)
	assert.NoError(t, err)
}
