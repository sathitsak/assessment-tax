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

	// "github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/sathitsak/assessment-tax/middleware"
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
	dbURL := os.Getenv("DATABASE_URL")
	adminID := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	db, err := db.New(dbURL)
	if err != nil {
		log.Fatal("can't connect to db")
	}
	e := echo.New()

	h := handler.CreateHandler(db)
	e.POST("/tax/calculations", middleware.ValidateRequestMiddleware(h.CalTaxHandler))
	e.POST("/tax/calculations/upload-csv", h.HandleFileUpload)
	g := e.Group("/admin")
	g.Use(middleware.ValidateBasicAuth(adminID, adminPassword))
	g.POST("/deductions/personal", h.PersonalAllowanceHandler)
	g.POST("/deductions/k-receipt", h.KReceiptHandler)


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



