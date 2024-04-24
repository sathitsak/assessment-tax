package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"


	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/sathitsak/assessment-tax/pkg/db"
	"github.com/sathitsak/assessment-tax/pkg/handler"
)

var PERSONAL_ALLOWANCE = 60000.0

func main() {
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	
	if err := db.Prepare(); err != nil {
		log.Fatal("can't connect to db")
	}
	e := echo.New()
	db,err := db.New()
	if err != nil {
		log.Fatal("can't connect to db")
	}
	h := handler.CreateHandler(db)
	e.POST("/tax/calculations", h.CalTaxHandler)
	e.POST("/admin/deductions/personal",h.PersonalAllowanceHandler)

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

