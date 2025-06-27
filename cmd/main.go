package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/martinusiron/loan-service/configs"
	"github.com/martinusiron/loan-service/repository/postgres"
	"github.com/martinusiron/loan-service/usecase"

	httpDelivery "github.com/martinusiron/loan-service/delivery/http"

	_ "github.com/lib/pq"
)

func main() {
	cfg := configs.Load()

	connStr := cfg.DB_URL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("failed to connect to db:", err)
	}
	defer db.Close()

	loanRepo := postgres.NewLoanRepo(db)
	approvalRepo := postgres.NewApprovalRepo(db)
	investRepo := postgres.NewInvestmentRepo(db)
	usecase := usecase.NewLoanUsecase(loanRepo, approvalRepo, investRepo)

	router := httpDelivery.InitRouter(&httpDelivery.Handler{UC: usecase})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Starting server on port", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
