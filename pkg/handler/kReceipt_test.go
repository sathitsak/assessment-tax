package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sathitsak/assessment-tax/internal"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	amount string
}

func TestSetInvalidKReceipt(t *testing.T) {
	tests := []string{
		`{"amount": 100000.1,}`,
		`{"amount": -0.1,}`,
	}
	for _, test := range tests {

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(test))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		db, teardown := internal.SetupTestDB(t)
		defer teardown()
		h := CreateHandler(db)
		// Assertions
		if assert.NoError(t, h.KReceiptHandler(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	}

}

func TestSetKReceipt(t *testing.T){
	var requestJSON = `{
		"amount": 70000.0
	  }`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	db,teardown:= internal.SetupTestDB(t)
	defer teardown()
	h:= CreateHandler(db)
	 want := KReceiptResponse{
		KReceipt: 70000.0,
	 }
	 var got KReceiptResponse
	if assert.NoError(t, h.KReceiptHandler(c),json.Unmarshal(rec.Body.Bytes(), &got)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t,got,want)
	}
	amount,err := h.kReceipt.Read()
	if assert.NoError(t,err){
		assert.Equal(t,70000.0,amount)

	}
	

}