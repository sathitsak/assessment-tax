package db

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)


func Prepare()error {
	db,err := New();
	if err != nil {
		return err
    }
	defer db.Close()
    query := `
    CREATE TABLE IF NOT EXISTS  personal_allowance (
		id serial PRIMARY KEY,
		amount double precision NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);`
	
    if _, err = db.Exec(query); err != nil{
			log.Fatalf("Error creating table: %s", err)
			return err
	}
	return nil
}

func New()(*sql.DB,error){
	err := godotenv.Load("../../.env")
	if err != nil {
	  log.Fatal("Error loading .env file")
	}
	  
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	
	if err != nil {
		log.Fatal(err)
		return nil,err
	}
	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return nil,err
	}
	return db,nil

}



