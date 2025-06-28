package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/martinusiron/loan-service/configs"
	"github.com/martinusiron/loan-service/delivery/http"
	"github.com/martinusiron/loan-service/repository/postgres"
	"github.com/martinusiron/loan-service/usecase"

	_ "github.com/lib/pq"
	_ "github.com/martinusiron/loan-service/docs" // Swagger docs
)

// @title Amartha Loan Service API
// @version 1.0
// @description RESTful API to simulate loan lifecycle (proposed → approved → invested → disbursed).
// @contact.name Martinus Iron Sijabat
// @contact.email your@email.com
// @host localhost:8080
// @BasePath /
func main() {
	cfg := configs.Load()

	db, err := sql.Open("postgres", cfg.DB_URL)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	loanRepo := postgres.NewLoanRepo(db)
	approvalRepo := postgres.NewApprovalRepo(db)
	investRepo := postgres.NewInvestmentRepo(db)

	uc := usecase.NewLoanUsecase(loanRepo, approvalRepo, investRepo, db)

	router := http.InitRouter(uc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Starting server on port %s...\n", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("server failed:", err)
	}
}
