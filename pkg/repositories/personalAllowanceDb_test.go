package repositories

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"
	"github.com/sathitsak/assessment-tax/pkg/db"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestSetPersonalAllowance(t *testing.T) {
	db, teardown := setup(t)
	want := 75000.0
	CreatePersonalAllowance(db, want)
	got, err := ReadPersonalAllowance(db)
	if assert.NoError(t, err) {
		assert.Equal(t, got, want)

	}
	teardown()
}

func setup(t *testing.T) (*sql.DB, func()) {
	err := godotenv.Load()
	assert.Equal(t, err, nil)
	db, err := db.New()
	assert.Equal(t, err, nil)
	
	setupTestTable(db)
	teardown := func() {
	dropTestTable(db)
		
	}
	return db, teardown
}



func setupTestTable(db *sql.DB) {
	// Create tables as required for your tests
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS personal_allowance (
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




