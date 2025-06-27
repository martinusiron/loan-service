package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/martinusiron/loan-service/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	UC *usecase.LoanUsecase
}

func NewHandler(r *gin.Engine, uc *usecase.LoanUsecase) {
	h := &Handler{UC: uc}

	v1 := r.Group("/v1")
	{
		v1.POST("/loans", h.CreateLoan)
		v1.POST("/loans/:id/approve", h.ApproveLoan)
		v1.POST("/loans/:id/invest", h.InvestLoan)
		v1.POST("/loans/:id/disburse", h.DisburseLoan)
		v1.GET("/loans/:id", h.GetLoan)
	}
}

type CreateLoanPayload struct {
	BorrowerID      string  `json:"borrower_id"`
	PrincipalAmount float64 `json:"principal_amount"`
	Rate            float64 `json:"rate"`
	ROI             float64 `json:"roi"`
}

type ApproveLoanPayload struct {
	PictureProof string `json:"picture_proof"`
	EmployeeID   string `json:"employee_id"`
	Date         string `json:"date"`
}

type InvestLoanPayload struct {
	InvestorEmail string  `json:"investor_email"`
	Amount        float64 `json:"amount"`
}

type DisburseLoanPayload struct {
	AgreementLink string `json:"agreement_letter_link"`
	EmployeeID    string `json:"employee_id"`
	Date          string `json:"date"`
}

// @Summary Create a new loan
// @Tags Loans
// @Accept json
// @Produce json
// @Param payload body CreateLoanPayload true "Loan payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/loans [post]
func (h *Handler) CreateLoan(c *gin.Context) {
	var payload CreateLoanPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loan, err := h.UC.CreateLoan(c.Request.Context(), payload.BorrowerID, payload.PrincipalAmount, payload.Rate, payload.ROI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, loan)
}

// @Summary Approve a loan
// @Tags Loans
// @Accept json
// @Produce json
// @Param id path int true "Loan ID"
// @Param payload body ApproveLoanPayload true "Approval payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /v1/loans/{id}/approve [post]
func (h *Handler) ApproveLoan(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var payload ApproveLoanPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	parsedDate, _ := time.Parse("2006-01-02", payload.Date)
	err := h.UC.ApproveLoan(c, id, payload.PictureProof, payload.EmployeeID, parsedDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Loan approved"})
}

// @Summary Invest in a loan
// @Tags Loans
// @Accept json
// @Produce json
// @Param id path int true "Loan ID"
// @Param payload body InvestLoanPayload true "Investment payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /v1/loans/{id}/invest [post]
func (h *Handler) InvestLoan(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var payload InvestLoanPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.UC.InvestLoan(c, id, payload.InvestorEmail, payload.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Investment accepted"})
}

// @Summary Disburse a loan
// @Tags Loans
// @Accept json
// @Produce json
// @Param id path int true "Loan ID"
// @Param payload body DisburseLoanPayload true "Disbursement payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /v1/loans/{id}/disburse [post]
func (h *Handler) DisburseLoan(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var payload DisburseLoanPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	parsedDate, _ := time.Parse("2006-01-02", payload.Date)
	err := h.UC.DisburseLoan(c, id, payload.AgreementLink, payload.EmployeeID, parsedDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Loan disbursed"})
}

// @Summary Get a loan by ID
// @Tags Loans
// @Produce json
// @Param id path int true "Loan ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /v1/loans/{id} [get]
func (h *Handler) GetLoan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	fmt.Println("id handler ", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error parsing id"})
		return
	}

	loan, err := h.UC.GetLoan(c, id)
	if err != nil || loan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "loan not found"})
		return
	}
	c.JSON(http.StatusOK, loan)
}
