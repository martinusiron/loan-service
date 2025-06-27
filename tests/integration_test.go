package tests

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/martinusiron/loan-service/delivery/http"
	"github.com/martinusiron/loan-service/repository/postgres"
	"github.com/martinusiron/loan-service/usecase"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	DB     *sql.DB
	Server *gin.Engine
}

func (s *IntegrationTestSuite) SetupSuite() {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "loan_db"),
	)

	db, err := sql.Open("postgres", dsn)
	s.Require().NoError(err)
	s.Require().NoError(db.Ping())

	s.DB = db

	loanRepo := postgres.NewLoanRepo(s.DB)
	approvalRepo := postgres.NewApprovalRepo(s.DB)
	investmentRepo := postgres.NewInvestmentRepo(s.DB)

	uc := usecase.NewLoanUsecase(loanRepo, approvalRepo, investmentRepo)

	r := gin.Default()
	http.NewHandler(r, uc)

	s.Server = r
}

func (s *IntegrationTestSuite) TearDownSuite() {
	if s.DB != nil {
		s.DB.Close()
	}
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

// helper
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
