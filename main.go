package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"fmt"
)

func main() {
	port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
        dbURL = "localhost:5432"
    }
	adminUsername := os.Getenv("ADMIN_USERNAME")
	if adminUsername == "" {
        adminUsername = "adminTax"
    }
	adminPassword := os.Getenv("ADMIN_USERNAME")
	if adminPassword == "" {
        adminPassword = "admin!"
    }
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
