package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestMiddlewareWithValidBody tests the middleware with a valid request body
func TestMiddlewareWithValidBody(t *testing.T) {
	e := echo.New()
	reqBody := map[string]interface{}{
		"totalIncome": 500000.0,
		"wht":         0.0,
		"allowances": []map[string]interface{}{
			{
				"allowanceType": "k-receipt",
				"amount":        200000.0,
			},
			{
				"allowanceType": "donation",
				"amount":        100000.0,
			},
		},
	}
	bodyData, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(bodyData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/test")

	handler := ValidateRequestMiddleware(func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"message": "Passed"})
	})

	if assert.NoError(t, handler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Passed")
	}
}

type testCase struct {
	reqBody     map[string]interface{}
	wantMessage string
}

func TestMiddlewareWithNoAllowances(t *testing.T) {
	e := echo.New()
	tests := []testCase{
		{reqBody: map[string]interface{}{
			"totalIncome": 500000.0,
			"wht":         0.0,
		}, wantMessage: "allowances are required"},
		{reqBody: map[string]interface{}{
			"totalIncome": 500000.0,
			"allowances": []map[string]interface{}{
				{
					"allowanceType": "k-receipt",
					"amount":        200000.0,
				},
				{
					"allowanceType": "donation",
					"amount":        100000.0,
				},
			},
		}, wantMessage: "wht is required"},
		{reqBody: map[string]interface{}{
			"wht": 0.0,
			"allowances": []map[string]interface{}{
				{
					"allowanceType": "k-receipt",
					"amount":        200000.0,
				},
				{
					"allowanceType": "donation",
					"amount":        100000.0,
				},
			},
		}, wantMessage: "totalIncome is required"},
		{reqBody: map[string]interface{}{
			"wht": 0.0,
			"allowances": []map[string]interface{}{
				{
					"allowanceType": "unknown allowanceType",
					"amount":        100000.0,
				},
			},
		}, wantMessage: "totalIncome is required"},
	}
	for _, test := range tests {

		bodyData, _ := json.Marshal(test.reqBody)
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(bodyData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/test")

		handler := ValidateRequestMiddleware(func(c echo.Context) error {
			return c.JSON(http.StatusOK, echo.Map{"message": "Passed"})
		})

		if assert.NoError(t, handler(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), test.wantMessage)
		}
	}
}
