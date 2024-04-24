package main

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sathitsak/assessment-tax/internal"
	"github.com/sathitsak/assessment-tax/middleware"
	"github.com/sathitsak/assessment-tax/pkg/handler"
	"github.com/stretchr/testify/assert"
)

func TestAdminWrongCredential(t *testing.T) {
	var requestJSON = `{}`
	req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", strings.NewReader(requestJSON))
	auth := base64.StdEncoding.EncodeToString([]byte("adminTaxx:adminn!"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Basic "+auth)
	rec := httptest.NewRecorder()
	db, teardown := internal.SetupTestDB(t)
	defer teardown()
	h := handler.CreateHandler(db)
	e := echo.New()
	e.Use(middleware.ValidateBasicAuth("adminTax", "admin!"))
	e.POST("/admin/deductions/personal", h.PersonalAllowanceHandler)
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}