package handler

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sathitsak/assessment-tax/internal"
	"github.com/stretchr/testify/assert"
)

func TestCalTaxHandler(t *testing.T) {
	var requestJSON = `{
		"totalIncome": 500000.0,
		"wht": 0.0,
		"allowances": [
			{
				"allowanceType": "k-receipt",
				"amount": 200000.0
			  },
			  {
				"allowanceType": "donation",
				"amount": 100000.0
			  }
		]
	  }`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	want := Response{
		Tax: 14000.0,
		TaxLevel: []TaxLevel{
			{
				Level: "0-150,000",
				Tax:   0.0,
			},
			{
				Level: "150,001-500,000",
				Tax:   14000.0,
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
	db,teardown:= internal.SetupTestDB(t)
	defer teardown()
	h:= CreateHandler(db)

	if assert.NoError(t, h.CalTaxHandler(c),json.Unmarshal(rec.Body.Bytes(), &got)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, want, got)
	}

}


func TestPersonalAllowanceHandler(t *testing.T) {
	var requestJSON = `{
		"amount": 70000.0
	  }`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", strings.NewReader(requestJSON))
	auth := base64.StdEncoding.EncodeToString([]byte("adminTax:admin!"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Basic "+auth)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	db,teardown:= internal.SetupTestDB(t)
	defer teardown()
	h:= CreateHandler(db)
	 want := PersonalAllowanceResponse{
		PersonalDeduction: 70000.0,
	 }
	 var got PersonalAllowanceResponse
	if assert.NoError(t, h.PersonalAllowanceHandler(c),json.Unmarshal(rec.Body.Bytes(), &got)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t,got,want)
	}
	amount,err := h.personalAllowance.Read()
	if assert.NoError(t,err){
		assert.Equal(t,70000.0,amount)

	}
	

}



func TestPerosnalAllowanceBadRequest(t *testing.T) {
	var requestJSON = `{}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", strings.NewReader(requestJSON))
	auth := base64.StdEncoding.EncodeToString([]byte("adminTax:admin!"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Basic "+auth)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	db,teardown:= internal.SetupTestDB(t)
	defer teardown()
	h:= CreateHandler(db)
	if assert.NoError(t, h.PersonalAllowanceHandler(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
	
}
func TestBadRequest(t *testing.T) {

	var badRequestJSON = `{"totalIncome": 500000.0,}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations", strings.NewReader(badRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	db, teardown := internal.SetupTestDB(t)
	h := CreateHandler(db)
	// Assertions
	if assert.NoError(t, h.CalTaxHandler(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
	teardown()
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
	db,teardown:= internal.SetupTestDB(t)
	defer teardown()
	h:= CreateHandler(db)

	if assert.NoError(t, h.CalTaxHandler(c),json.Unmarshal(rec.Body.Bytes(), &got)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, want, float64(got.TaxRefund))
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
	db,teardown:= internal.SetupTestDB(t)
	defer teardown()
	h:= CreateHandler(db)

	if assert.NoError(t, h.CalTaxHandler(c),json.Unmarshal(rec.Body.Bytes(), &got)) {
		assert.Equal(t, want, float64(got.Tax))
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
