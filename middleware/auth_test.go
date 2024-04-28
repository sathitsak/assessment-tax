package middleware

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sathitsak/assessment-tax/pkg/handler"
	"github.com/stretchr/testify/assert"
)

// TestValidateBasicAuth tests the Basic Auth middleware for both valid and invalid credentials.
func TestValidateBasicAuth(t *testing.T) {
	// Set the expected username and password
	expectedUsername := "adminTax"
	expectedPassword := "admin!"

	// Create a new instance of Echo
	e := echo.New()

	// Create a test handler that will be protected by the middleware
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test passed")
	}

	// Tests
	test := struct {
		name           string
		username       string
		password       string
		expectedStatus int
	}{

		name:           "Valid credentials",
		username:       "adminTax",
		password:       "admin!",
		expectedStatus: http.StatusOK,
	}

	// Run tests

	t.Run(test.name, func(t *testing.T) {
		// Create a request
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetBasicAuth(test.username, test.password)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Attach the middleware to the handler
		mw := ValidateBasicAuth(expectedUsername, expectedPassword)
		handlerWithMiddleware := mw(handler)

		// Execute the handler
		if assert.NoError(t, handlerWithMiddleware(c)) {
			assert.Equal(t, test.expectedStatus, rec.Code)
		}
	})

}

func TestInValidateBasicAuth(t *testing.T) {
	// Test cases
	testCases := []struct {
		name               string
		username           string
		password           string
		expectedUsername   string
		expectedPassword   string
		expectedStatusCode int
	}{
		{
			name:               "Valid credentials",
			username:           "validuser",
			password:           "validpassword",
			expectedUsername:   "validuser",
			expectedPassword:   "validpassword",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid username",
			username:           "invaliduser",
			password:           "validpassword",
			expectedUsername:   "validuser",
			expectedPassword:   "validpassword",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "Invalid password",
			username:           "validuser",
			password:           "invalidpassword",
			expectedUsername:   "validuser",
			expectedPassword:   "validpassword",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			var requestJSON = `{}`
			req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", strings.NewReader(requestJSON))
			auth := base64.StdEncoding.EncodeToString([]byte("adminTaxx:adminn!"))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set("Authorization", "Basic "+auth)
			rec := httptest.NewRecorder()

			h := handler.CreateTestHandler()
			e := echo.New()
			e.Use(ValidateBasicAuth("adminTax", "admin!"))
			e.POST("/admin/deductions/personal", h.PersonalAllowanceHandler)
			e.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusUnauthorized, rec.Code)

		})
	}
}
