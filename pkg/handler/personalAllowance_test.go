package handler



import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)


func TestSetInvalidPersonalAllowance(t *testing.T) {
	tests := []string{
		`{"amount": 100000.1,}`,
		`{"amount": 9999.9,}`,
	}
	for _, test := range tests {

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(test))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		h := CreateTestHandler()
		// Assertions
		if assert.NoError(t, h.PersonalAllowanceHandler(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	}

}

func TestSetPersonalAllowance(t *testing.T) {
	var requestJSON = `{
		"amount": 70000.0
	  }`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := CreateTestHandler()
	want := PersonalAllowanceResponse{
		PersonalDeduction: 70000.0,
	}
	var got PersonalAllowanceResponse
	if assert.NoError(t, h.PersonalAllowanceHandler(c), json.Unmarshal(rec.Body.Bytes(), &got)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, got, want)
	}
	amount, err := h.personalAllowance.Read()
	if assert.NoError(t, err) {
		assert.Equal(t, 70000.0, amount)

	}

}
