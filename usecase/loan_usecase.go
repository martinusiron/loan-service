package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/martinusiron/loan-service/domain"
	"github.com/martinusiron/loan-service/repository"
	"github.com/martinusiron/loan-service/utils"
)

type LoanUsecase struct {
	LoanRepo       repository.LoanRepository
	ApprovalRepo   repository.ApprovalRepository
	InvestmentRepo repository.InvestmentRepository
}

func NewLoanUsecase(loan repository.LoanRepository, approval repository.ApprovalRepository, invest repository.InvestmentRepository) *LoanUsecase {
	return &LoanUsecase{
		LoanRepo:       loan,
		ApprovalRepo:   approval,
		InvestmentRepo: invest,
	}
}

func (uc *LoanUsecase) CreateLoan(ctx context.Context, borrowerID string, principal, rate, roi float64) (*domain.Loan, error) {
	loan := &domain.Loan{
		BorrowerID:      borrowerID,
		PrincipalAmount: principal,
		Rate:            rate,
		ROI:             roi,
		Status:          domain.StatusProposed,
	}
	err := uc.LoanRepo.CreateLoan(ctx, loan)
	return loan, err
}

func (uc *LoanUsecase) ApproveLoan(ctx context.Context, loanID int, pictureProof, employeeID string, date time.Time) error {
	loan, err := uc.LoanRepo.GetLoanByID(ctx, loanID)
	if err != nil || loan == nil {
		return errors.New("loan not found")
	}
	if loan.Status != domain.StatusProposed {
		return errors.New("loan is not in proposed state")
	}
	approval := &domain.LoanApproval{
		LoanID:       loanID,
		PictureProof: pictureProof,
		EmployeeID:   employeeID,
		ApprovedAt:   date,
	}
	if err := uc.ApprovalRepo.CreateApproval(ctx, approval); err != nil {
		return err
	}
	return uc.LoanRepo.UpdateLoanStatus(ctx, loanID, domain.StatusApproved)
}

func (uc *LoanUsecase) InvestLoan(ctx context.Context, loanID int, investorEmail string, amount float64) error {
	loan, err := uc.LoanRepo.GetLoanByID(ctx, loanID)
	if err != nil || loan == nil {
		return errors.New("loan not found")
	}
	if loan.Status != domain.StatusApproved && loan.Status != domain.StatusInvested {
		return errors.New("loan not available for investment")
	}
	totalInvested, err := uc.InvestmentRepo.GetTotalInvested(ctx, loanID)
	if err != nil {
		return err
	}
	if totalInvested+amount > loan.PrincipalAmount {
		return fmt.Errorf("investment exceeds loan principal")
	}
	inv := &domain.Investment{
		LoanID:        loanID,
		InvestorEmail: investorEmail,
		Amount:        amount,
	}
	if err := uc.InvestmentRepo.AddInvestment(ctx, inv); err != nil {
		return err
	}

	// Check again if fully funded
	newTotal := totalInvested + amount
	if newTotal == loan.PrincipalAmount && loan.Status != domain.StatusInvested {
		if err := uc.LoanRepo.UpdateLoanStatus(ctx, loanID, domain.StatusInvested); err != nil {
			return err
		}
		// Email notification
		investors, _ := uc.InvestmentRepo.GetInvestorsByLoan(ctx, loanID)
		for _, i := range investors {
			utils.SendEmail(i.InvestorEmail, fmt.Sprintf("Loan #%d funded", loanID), "Agreement letter link will follow.")
		}
	}
	return nil
}

func (uc *LoanUsecase) DisburseLoan(ctx context.Context, loanID int, agreementLink, employeeID string, date time.Time) error {
	loan, err := uc.LoanRepo.GetLoanByID(ctx, loanID)
	if err != nil || loan == nil {
		return errors.New("loan not found")
	}
	if loan.Status != domain.StatusInvested {
		return errors.New("loan is not ready for disbursement")
	}
	if err := uc.LoanRepo.SetAgreementLink(ctx, loanID, agreementLink); err != nil {
		return err
	}
	return uc.LoanRepo.UpdateLoanStatus(ctx, loanID, domain.StatusDisbursed)
}

func (uc *LoanUsecase) GetLoan(ctx context.Context, loanID int) (*domain.Loan, error) {
	loan, err := uc.LoanRepo.GetLoanByID(ctx, loanID)
	if err != nil {
		fmt.Println("usecase ", err)
		return nil, err
	}

	return loan, nil
}
