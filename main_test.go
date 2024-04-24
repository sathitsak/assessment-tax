

package main

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sathitsak/assessment-tax/pkg/handler"
	"github.com/stretchr/testify/assert"
	"github.com/sathitsak/assessment-tax/internal"

)

func TestBadRequest(t *testing.T){
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}
	  dbURL := os.Getenv("DATABASE_URL")
	  db, err := sql.Open("postgres", dbURL)
	  if err != nil {
		  log.Fatal(err)
	  }
	  defer db.Close()
  
	  // Check the connection
	  err = db.Ping()
	  if err != nil {
		  log.Fatal(err)
	  }
	tests := []string{`{"amount": 100000.1,}`,`{"amount": 9999.9,}`} 
	for _,badRequestJSON := range tests{
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", strings.NewReader(badRequestJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		db, teardown := internal.SetupTestDB(t)
		h := handler.CreateHandler(db)
		// Assertions
		if assert.NoError(t, h.PersonalAllowanceHandler(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
		teardown()
	}
}

