package postgres

import (
	"context"
	"database/sql"

	"github.com/martinusiron/loan-service/domain"
	"github.com/martinusiron/loan-service/utils"
)

type ApprovalRepo struct {
	DB *sql.DB
}

func NewApprovalRepo(db *sql.DB) *ApprovalRepo {
	return &ApprovalRepo{DB: db}
}

func (r *ApprovalRepo) CreateApproval(ctx context.Context, a *domain.LoanApproval) error {
	exec := utils.GetExecutor(ctx, r.DB)
	query := `INSERT INTO loan_approvals (loan_id, picture_proof, employee_id, approved_at) VALUES ($1, $2, $3, $4)`

	_, err := exec.ExecContext(ctx, query, a.LoanID, a.PictureProof, a.EmployeeID, a.ApprovedAt)
	return err
}
