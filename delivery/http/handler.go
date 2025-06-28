package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/martinusiron/loan-service/dto"
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

func errorResponse(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
}

// @Summary Create a new loan
// @Tags Loans
// @Accept json
// @Produce json
// @Param payload body dto.CreateLoanPayload true "Loan payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/loans [post]
func (h *Handler) CreateLoan(c *gin.Context) {
	var payload dto.CreateLoanPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	loan, err := h.UC.CreateLoan(
		c.Request.Context(),
		payload,
	)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, loan)
}

// @Summary Approve a loan
// @Tags Loans
// @Accept json
// @Produce json
// @Param id path int true "Loan ID"
// @Param payload body dto.ApproveLoanPayload true "Approval payload (date format: YYYY-MM-DD)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /v1/loans/{id}/approve [post]
func (h *Handler) ApproveLoan(c *gin.Context) {
	var payload dto.ApproveLoanPayload

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	payload.LoanID = id
	if err := c.ShouldBindJSON(&payload); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	parsedDate, err := time.Parse("2006-01-02", payload.DateStr)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid date format, must be YYYY-MM-DD"))
		return
	}
	payload.Date = parsedDate

	if err := h.UC.ApproveLoan(c, payload); err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan approved"})
}

// @Summary Invest in a loan
// @Tags Loans
// @Accept json
// @Produce json
// @Param id path int true "Loan ID"
// @Param payload body dto.InvestLoanPayload true "Investment payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /v1/loans/{id}/invest [post]
func (h *Handler) InvestLoan(c *gin.Context) {
	var payload dto.InvestLoanPayload

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	payload.LoanID = id
	if err := c.ShouldBindJSON(&payload); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := h.UC.InvestLoan(c,
		payload,
	); err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Investment accepted"})
}

// @Summary Disburse a loan
// @Tags Loans
// @Accept json
// @Produce json
// @Param id path int true "Loan ID"
// @Param payload body dto.DisburseLoanPayload true "Disbursement payload (date format: YYYY-MM-DD)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /v1/loans/{id}/disburse [post]
func (h *Handler) DisburseLoan(c *gin.Context) {
	var payload dto.DisburseLoanPayload

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	payload.LoanID = id

	if err := c.ShouldBindJSON(&payload); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	parsedDate, err := time.Parse("2006-01-02", payload.DateStr)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid date format, must be YYYY-MM-DD"))
		return
	}
	payload.Date = parsedDate

	if err := h.UC.DisburseLoan(c,
		payload,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	loan, err := h.UC.GetLoan(c, id)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if loan == nil {
		errorResponse(c, http.StatusNotFound, errors.New("loan not found"))
		return
	}

	c.JSON(http.StatusOK, loan)
}
