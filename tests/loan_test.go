package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

func (s *IntegrationTestSuite) TestCreateAndGetLoan() {
	payload := map[string]interface{}{
		"borrower_id":      "BR01",
		"principal_amount": 1000000,
		"rate":             10.0,
		"roi":              5.0,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/v1/loans", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Server.ServeHTTP(w, req)

	s.Equal(201, w.Code)
	s.T().Log("CreateLoan response:", w.Body.String())

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err)
	idVal, ok := resp["ID"].(float64)
	s.Require().True(ok, "Expected 'ID' in response")
	loanID := int(idVal)

	req2 := httptest.NewRequest(http.MethodGet, "/v1/loans/"+itoa(loanID), nil)
	w2 := httptest.NewRecorder()
	s.Server.ServeHTTP(w2, req2)

	s.Equal(200, w2.Code)
	s.T().Log("GetLoan response:", w2.Body.String())
}

func (s *IntegrationTestSuite) TestApproveLoan() {
	create := map[string]interface{}{
		"borrower_id":      "BR02",
		"principal_amount": 500000,
		"rate":             12.0,
		"roi":              8.0,
	}
	body, _ := json.Marshal(create)
	req := httptest.NewRequest(http.MethodPost, "/v1/loans", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Server.ServeHTTP(w, req)
	s.Equal(201, w.Code)
	s.T().Log("CreateLoan response:", w.Body.String())

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err)
	idVal, ok := resp["ID"].(float64)
	s.Require().True(ok, "Expected 'ID' in response")
	loanID := int(idVal)

	approve := map[string]interface{}{
		"loan_id":       loanID,
		"picture_proof": "proof.jpg",
		"employee_id":   "EMP001",
		"date":          "2025-06-26",
	}
	apBody, _ := json.Marshal(approve)
	req2 := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/v1/loans/%d/approve", loanID), bytes.NewBuffer(apBody))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	s.Server.ServeHTTP(w2, req2)

	s.Equal(200, w2.Code)
	s.T().Log("ApproveLoan response:", w2.Body.String())
}

func (s *IntegrationTestSuite) TestInvestLoanAndDisburse() {
	create := map[string]interface{}{
		"borrower_id":      "BR03",
		"principal_amount": 1000000,
		"rate":             10.0,
		"roi":              5.0,
	}
	body, _ := json.Marshal(create)
	req := httptest.NewRequest(http.MethodPost, "/v1/loans", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Server.ServeHTTP(w, req)
	s.Equal(201, w.Code)
	s.T().Log("CreateLoan response:", w.Body.String())

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err)
	idVal, ok := resp["ID"].(float64)
	s.Require().True(ok, "Expected 'ID' in response")
	loanID := int(idVal)

	ap := map[string]interface{}{
		"loan_id":       loanID,
		"picture_proof": "proof.jpg",
		"employee_id":   "EMP002",
		"date":          "2025-06-26",
	}
	apBody, _ := json.Marshal(ap)
	req2 := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/v1/loans/%d/approve", loanID), bytes.NewBuffer(apBody))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	s.Server.ServeHTTP(w2, req2)
	s.Equal(200, w2.Code)
	s.T().Log("ApproveLoan response:", w2.Body.String())

	invest := map[string]interface{}{
		"loan_id":        loanID,
		"investor_email": "foo@bar.com",
		"amount":         1000000,
	}
	invBody, _ := json.Marshal(invest)
	req3 := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/v1/loans/%d/invest", loanID), bytes.NewBuffer(invBody))
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	s.Server.ServeHTTP(w3, req3)
	s.Equal(200, w3.Code)
	s.T().Log("InvestLoan response:", w3.Body.String())

	dis := map[string]interface{}{
		"loan_id":               loanID,
		"agreement_letter_link": "http://example.com/agreement.pdf",
		"employee_id":           "EMP003",
		"date":                  "2025-06-26",
	}
	disBody, _ := json.Marshal(dis)
	req4 := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/v1/loans/%d/disburse", loanID), bytes.NewBuffer(disBody))
	req4.Header.Set("Content-Type", "application/json")
	w4 := httptest.NewRecorder()
	s.Server.ServeHTTP(w4, req4)
	s.Equal(200, w4.Code)
	s.T().Log("DisburseLoan response:", w4.Body.String())
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
