package internal

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func SetupTestDB(t *testing.T) (*sql.DB, func()) {
	
	db, err := createTestDB()
	assert.Equal(t, err, nil)

	setupTestTable(db)
	teardown := func() {
		dropTestTable(db)

	}
	return db, teardown
}

func setupTestTable(db *sql.DB) {
	// Create tables as required for your tests
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS  personal_allowance (
		id serial PRIMARY KEY,
		amount double precision NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

func dropTestTable(db *sql.DB) {
	_, err := db.Exec(`DROP TABLE IF EXISTS personal_allowance`)
	if err != nil {
		log.Printf("Failed to drop table: %v", err)
	}
}
func createTestDB() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return db, nil

}
