package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/martinusiron/loan-service/domain"
	"github.com/martinusiron/loan-service/utils"
)

type LoanRepo struct {
	DB *sql.DB
}

func NewLoanRepo(db *sql.DB) *LoanRepo {
	return &LoanRepo{DB: db}
}

func (r *LoanRepo) CreateLoan(ctx context.Context, loan *domain.Loan) error {
	exec := utils.GetExecutor(ctx, r.DB)
	query := `INSERT INTO loans (borrower_id, principal_amount, rate, roi) VALUES ($1, $2, $3, $4) RETURNING id`
	return exec.QueryRowContext(ctx, query,
		loan.BorrowerID, loan.PrincipalAmount, loan.Rate, loan.ROI).Scan(&loan.ID)
}

func (r *LoanRepo) GetLoanByID(ctx context.Context, id int) (*domain.Loan, error) {
	exec := utils.GetExecutor(ctx, r.DB)
	query := `SELECT id, borrower_id, principal_amount, rate, roi, status, COALESCE(agreement_letter_link, '') AS agreement_letter_link, created_at, updated_at FROM loans WHERE id = $1`
	row := exec.QueryRowContext(ctx, query, id)

	var l domain.Loan
	err := row.Scan(
		&l.ID,
		&l.BorrowerID,
		&l.PrincipalAmount,
		&l.Rate,
		&l.ROI,
		&l.Status,
		&l.AgreementLetterLink,
		&l.CreatedAt,
		&l.UpdatedAt,
	)
	if err != nil {
		fmt.Println("repo ", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &l, nil
}

func (r *LoanRepo) UpdateLoanStatus(ctx context.Context, id int, status domain.LoanStatus) error {
	exec := utils.GetExecutor(ctx, r.DB)
	_, err := exec.ExecContext(ctx, `UPDATE loans SET status = $1, updated_at = NOW() WHERE id = $2`, status, id)
	return err
}

func (r *LoanRepo) SetAgreementLink(ctx context.Context, id int, link string) error {
	exec := utils.GetExecutor(ctx, r.DB)
	_, err := exec.ExecContext(ctx, `UPDATE loans SET agreement_letter_link = $1, updated_at = NOW() WHERE id = $2`, link, id)
	return err
}
