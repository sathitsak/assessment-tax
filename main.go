package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"

	// "github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/sathitsak/assessment-tax/internal"
	"github.com/sathitsak/assessment-tax/internal/db"
	"github.com/sathitsak/assessment-tax/middleware"
	"github.com/sathitsak/assessment-tax/pkg/handler"
)

var PERSONAL_ALLOWANCE = 60000.0

func main() {
	port := internal.GetEnvWithFallback("PORT", "8080")
	dbURL := internal.GetEnvWithFallback("DATABASE_URL", "postgres://root:password@localhost:15432/tax_assessment?sslmode=disable")
	adminID := internal.GetEnvWithFallback("ADMIN_USERNAME", "adminTax")
	adminPassword := internal.GetEnvWithFallback("ADMIN_PASSWORD", "admin!")
	fmt.Println(dbURL)
	db, err := db.New(dbURL)
	if err != nil {
		log.Fatal("can't connect to db")
	}
	e := echo.New()

	h := handler.CreateHandler(db)
	e.POST("/tax/calculations", h.CalTaxHandler)

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
