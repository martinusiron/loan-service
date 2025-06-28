package postgres

import (
	"context"
	"database/sql"

	"github.com/martinusiron/loan-service/domain"
	"github.com/martinusiron/loan-service/utils"
)

type InvestmentRepo struct {
	DB *sql.DB
}

func NewInvestmentRepo(db *sql.DB) *InvestmentRepo {
	return &InvestmentRepo{DB: db}
}

func (r *InvestmentRepo) AddInvestment(ctx context.Context, i *domain.Investment) error {
	exec := utils.GetExecutor(ctx, r.DB)
	query := `INSERT INTO investments (loan_id, investor_email, amount) VALUES ($1, $2, $3)`

	_, err := exec.ExecContext(ctx, query, i.LoanID, i.InvestorEmail, i.Amount)
	return err
}

func (r *InvestmentRepo) GetTotalInvested(ctx context.Context, loanID int) (float64, error) {
	exec := utils.GetExecutor(ctx, r.DB)
	query := `SELECT COALESCE(SUM(amount), 0) FROM investments WHERE loan_id = $1`

	var total float64
	err := exec.QueryRowContext(ctx, query, loanID).Scan(&total)
	return total, err
}

func (r *InvestmentRepo) GetInvestorsByLoan(ctx context.Context, loanID int) ([]domain.Investment, error) {
	exec := utils.GetExecutor(ctx, r.DB)
	query := `SELECT id, loan_id, investor_email, amount, invested_at FROM investments WHERE loan_id = $1`

	rows, err := exec.QueryContext(ctx, query, loanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var investors []domain.Investment
	for rows.Next() {
		var i domain.Investment
		if err := rows.Scan(&i.ID, &i.LoanID, &i.InvestorEmail, &i.Amount, &i.InvestedAt); err != nil {
			return nil, err
		}
		investors = append(investors, i)
	}
	return investors, nil
}
