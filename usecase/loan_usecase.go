package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/martinusiron/loan-service/domain"
	"github.com/martinusiron/loan-service/dto"
	"github.com/martinusiron/loan-service/repository"
	"github.com/martinusiron/loan-service/utils"
)

type LoanUsecase struct {
	LoanRepo       repository.LoanRepository
	ApprovalRepo   repository.ApprovalRepository
	InvestmentRepo repository.InvestmentRepository
	DB             *sql.DB
}

func NewLoanUsecase(lr repository.LoanRepository, ar repository.ApprovalRepository, ir repository.InvestmentRepository, db *sql.DB) *LoanUsecase {
	return &LoanUsecase{
		LoanRepo:       lr,
		ApprovalRepo:   ar,
		InvestmentRepo: ir,
		DB:             db,
	}
}

func (uc *LoanUsecase) CreateLoan(ctx context.Context, payload dto.CreateLoanPayload) (*domain.Loan, error) {
	loan := &domain.Loan{
		BorrowerID:      payload.BorrowerID,
		PrincipalAmount: payload.PrincipalAmount,
		Rate:            payload.Rate,
		ROI:             payload.ROI,
		Status:          domain.StatusProposed,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := uc.LoanRepo.CreateLoan(ctx, loan); err != nil {
		return nil, err
	}

	return loan, nil
}

func (uc *LoanUsecase) ApproveLoan(ctx context.Context, payload dto.ApproveLoanPayload) error {
	return utils.WithTransaction(ctx, uc.DB, func(txCtx context.Context) error {
		loan, err := uc.LoanRepo.GetLoanByID(txCtx, payload.LoanID)
		if err != nil || loan == nil {
			return errors.New("loan not found")
		}

		if loan.Status != domain.StatusProposed {
			return errors.New("loan is not in proposed state")
		}

		if err := uc.ApprovalRepo.CreateApproval(txCtx, &domain.LoanApproval{
			LoanID:       payload.LoanID,
			PictureProof: payload.PictureProof,
			EmployeeID:   payload.EmployeeID,
			ApprovedAt:   payload.Date,
		}); err != nil {
			return err
		}

		return uc.LoanRepo.UpdateLoanStatus(txCtx, payload.LoanID, domain.StatusApproved)
	})
}

func (uc *LoanUsecase) InvestLoan(ctx context.Context, payload dto.InvestLoanPayload) error {
	return utils.WithTransaction(ctx, uc.DB, func(txCtx context.Context) error {
		loan, err := uc.LoanRepo.GetLoanByID(txCtx, payload.LoanID)
		if err != nil || loan == nil {
			return errors.New("loan not found")
		}

		if loan.Status != domain.StatusApproved && loan.Status != domain.StatusInvested {
			return errors.New("loan not available for investment")
		}

		totalInvested, err := uc.InvestmentRepo.GetTotalInvested(txCtx, payload.LoanID)
		if err != nil {
			return err
		}

		if totalInvested+payload.Amount > loan.PrincipalAmount {
			return fmt.Errorf("investment exceeds loan principal")
		}

		if err := uc.InvestmentRepo.AddInvestment(txCtx, &domain.Investment{
			LoanID:        payload.LoanID,
			InvestorEmail: payload.InvestorEmail,
			Amount:        payload.Amount,
			InvestedAt:    time.Now(),
		}); err != nil {
			return err
		}

		newTotal := totalInvested + payload.Amount
		if newTotal == loan.PrincipalAmount && loan.Status != domain.StatusInvested {
			if err := uc.LoanRepo.UpdateLoanStatus(txCtx, payload.LoanID, domain.StatusInvested); err != nil {
				return err
			}
			investors, _ := uc.InvestmentRepo.GetInvestorsByLoan(txCtx, payload.LoanID)
			for _, inv := range investors {
				go utils.SendDummyEmail(inv.InvestorEmail, loan.ID)
			}
		}
		return nil
	})
}

func (uc *LoanUsecase) DisburseLoan(ctx context.Context, payload dto.DisburseLoanPayload) error {
	return utils.WithTransaction(ctx, uc.DB, func(txCtx context.Context) error {
		loan, err := uc.LoanRepo.GetLoanByID(txCtx, payload.LoanID)
		if err != nil || loan == nil {
			return errors.New("loan not found")
		}

		if loan.Status != domain.StatusInvested {
			return errors.New("loan is not ready for disbursement")
		}

		if err := uc.LoanRepo.SetAgreementLink(txCtx, payload.LoanID, payload.AgreementLink); err != nil {
			return err
		}

		return uc.LoanRepo.UpdateLoanStatus(txCtx, payload.LoanID, domain.StatusDisbursed)
	})
}

func (uc *LoanUsecase) GetLoan(ctx context.Context, id int) (*domain.Loan, error) {
	
	return uc.LoanRepo.GetLoanByID(ctx, id)
}
