package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"database/sql"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/sathitsak/assessment-tax/pkg/handler"
)

var PERSONAL_ALLOWANCE = 60000.0

func main() {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
	port := os.Getenv("PORT")
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

    log.Println("Successfully connected!")
	s := Server{db:db}
	s.createTable()
	s.seedTable()
	e := echo.New()
	h := handler.CreateHandler()
	e.POST("/tax/calculations", h.CalTaxHandler)
	e.POST("/admin/deductions/personal",s.setPersonalAllowance)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
type Server struct {
    db *sql.DB
}

func (s *Server)createTable() {
    query := `
    CREATE TABLE IF NOT EXISTS  personal_allowance (
		id serial PRIMARY KEY,
		amount double precision NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);`
    _, err := s.db.Exec(query)
    if err != nil {
        log.Fatalf("Error creating table: %s", err)
    }
}

type PersonalAllowance struct{
	Amount float64 `json:"amount" form:"amount"`
}

func (s *Server) setPersonalAllowance (c echo.Context) error {
	var pa PersonalAllowance
	if err := c.Bind(&pa); err != nil{
		return c.String(http.StatusBadRequest, "bad request")
	}
	fmt.Println(pa.Amount)
	if pa.Amount > 100000.0  {
		return c.String(http.StatusBadRequest, "The amount provided exceeds the maximum allowed limit.")
	}
	if pa.Amount < 10000.0 {
		return c.String(http.StatusBadRequest, "The amount provided is below the minimum allowed limit.")
	}
	query := fmt.Sprintf("INSERT INTO personal_allowance (amount) VALUES (%f);",pa.Amount)
	_,err := s.db.Exec(query)
	if  err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error please contact admin or try again later")
	}
	return c.JSON(http.StatusAccepted,pa)

}
func (s *Server)seedTable(){
	query := `SELECT amount FROM personal_allowance  ORDER BY created_at DESC LIMIT 1;`
	// Declare a variable to store the data from the row.
    var amount float64
    

    // Execute the query.
    row := s.db.QueryRow(query)
    err := row.Scan(&amount)
    if err != nil {
        if err == sql.ErrNoRows {
			s.db.Exec(`INSERT INTO personal_allowance (amount)
			VALUES (60000.0); `)
            fmt.Println("No rows were returned!")
        } else {
            log.Fatal(err)
        }
    } else {
        fmt.Printf("current amount: %f", amount)
    }
}

