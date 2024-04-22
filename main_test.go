// +build integration
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
	"github.com/stretchr/testify/assert"
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
		s := Server{db:db}
		// Assertions
		if assert.NoError(t, s.setPersonalAllowance(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	}
}

func TestSetPersonalAllowance(t *testing.T){
	db, teardown := setup(t)
	s := Server{db:db}
	want:= 75000.0

	_,err2:= s.InsertPersonalAllowance(want)
	got,err := s.ReadPersonalAllowance()
	if assert.NoError(t,err,err2) {
		assert.Equal(t, got,want)

	}
	teardown()
}


func setup(t *testing.T) (*sql.DB, func()) {
	err := godotenv.Load()
	assert.Equal(t, err, nil)
	dbURL := os.Getenv("DATABASE_URL")
	db,err := NewDB(dbURL)
	assert.Equal(t, err, nil)
	dropTestDB(db)
	createTestDB(db)
	setupTestTable(db)
	teardown := func() {
		dropTestTable(db)
		dropTestDB(db)
	}
	return db, teardown
}

func createTestDB(db *sql.DB) {
    _, err := db.Exec(`CREATE DATABASE test_db`)
    if err != nil {
        log.Fatalf("Failed to create test database: %v", err)
    }
}

func setupTestTable(db *sql.DB) {
    // Create tables as required for your tests
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
			id serial PRIMARY KEY,
			amount double precision NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
        );
    `)
    if err != nil {
        log.Fatalf("Failed to create table: %v", err)
    }
}

func dropTestTable(db *sql.DB) {
    _, err := db.Exec(`DROP TABLE IF EXISTS users`)
    if err != nil {
        log.Printf("Failed to drop table: %v", err)
    }
}

func dropTestDB(db *sql.DB) {
    _, err := db.Exec(`DROP DATABASE IF EXISTS test_db`)
    if err != nil {
        log.Printf("Failed to drop test database: %v", err)
    }
}
