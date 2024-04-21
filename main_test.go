package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)


  

func TestCalTaxHandler(t *testing.T) {
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
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	want := `{"tax":29000,"taxRefund":0}`

	c := e.NewContext(req, rec)
	
	if assert.NoError(t, calTaxHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, want, strings.TrimSpace(rec.Body.String()))
	}
}

func TestBadRequest(t *testing.T) {

	var badRequestJSON = `{"totalIncome": 500000.0,}`
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

func TestTaxRefund(t *testing.T){
	var requestJSON = `{
		"totalIncome": 500000.0,
		"wht": 39000.0,
		"allowances": [
		  {
			"allowanceType": "donation",
			"amount": 0.0
		  }
		]
	  }`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	want := `{"tax":0,"taxRefund":10000}`

	c := e.NewContext(req, rec)
	
	if assert.NoError(t, calTaxHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, want, strings.TrimSpace(rec.Body.String()))
	}
}

func TestDonation(t *testing.T){
	var requestJSON = `{
		"totalIncome": 500000.0,
		"wht": 0.0,
		"allowances": [
		  {
			"allowanceType": "donation",
			"amount": 200000.0
		  }
		]
	  }`
	  e := echo.New()
	  req := httptest.NewRequest(http.MethodPost, "/tax/calculations", strings.NewReader(requestJSON))
	  req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	  rec := httptest.NewRecorder()
	  want := `{"tax":19000,"taxRefund":0}`
  
	  c := e.NewContext(req, rec)
	  
	  if assert.NoError(t, calTaxHandler(c)) {
		  assert.Equal(t, http.StatusOK, rec.Code)
		  assert.Equal(t, want, strings.TrimSpace(rec.Body.String()))
	  }  
}
