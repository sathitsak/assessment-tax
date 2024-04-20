package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var requestJSON = `{
	"totalIncome": 500000.0,
	"wht": 0.0,
	"allowances": [
	  {
		"allowanceType": "donation",
		"amount": 0.0
	  }
	]
  }`
  var badRequestJSON = `{
	"totalIncome": 500000.0,
	
  }`

func TestCalTaxHandler(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	resJSON := `{"tax":29000}`
	

	// Assertions
	if assert.NoError(t, calTaxHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, resJSON, strings.TrimSpace(rec.Body.String()))
	}
}

func TestBadRequest(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations", strings.NewReader(badRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	

	// Assertions
	if assert.NoError(t, calTaxHandler(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

