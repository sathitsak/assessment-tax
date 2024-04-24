package db

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
)



func New(dbURL string)(*sql.DB,error){
	
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
	query := `
    CREATE TABLE IF NOT EXISTS  personal_allowance (
		id serial PRIMARY KEY,
		amount double precision NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);`
	
    if _, err := db.Exec(query); err != nil{
			log.Fatalf("Error creating table: %s", err)
			return nil,err
	}
	return db,nil

}



