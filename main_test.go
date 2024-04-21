package main

import (
	"encoding/json"
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
			"amount": 200000.0
		  }
		]
	  }`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	want := Response{
		Tax: 19000.0,
		TaxLevel: []TaxLevel{
			{
				Level: "0-150,000",
				Tax:   0.0,
			},
			{
				Level: "150,001-500,000",
				Tax:   19000.0,
			},
			{
				Level: "500,001-1,000,000",
				Tax:   0.0,
			},
			{
				Level: "1,000,001-2,000,000",
				Tax:   0.0,
			},
			{
				Level: "2,000,001 ขึ้นไป",
				Tax:   0.0,
			},
		},
	}

	c := e.NewContext(req, rec)
	var got Response

	if assert.NoError(t, ResponseToJSON(c, rec, &got)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, want, got)
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

func TestTaxRefund(t *testing.T) {
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
	want := 10000.0

	c := e.NewContext(req, rec)
	var got Response

	if assert.NoError(t, ResponseToJSON(c, rec, &got)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, want, got.TaxRefund)
	}
}

func TestDonation(t *testing.T) {
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
	want := 19000.0
	c := e.NewContext(req, rec)
	var got Response

	if assert.NoError(t, ResponseToJSON(c, rec, &got)) {
		assert.Equal(t, want, got.Tax)
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func ResponseToJSON(c echo.Context, rec *httptest.ResponseRecorder, data *Response) error {
	if err := calTaxHandler(c); err != nil {
		return err
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &data); err != nil {
		return err
	}
	return nil
}
