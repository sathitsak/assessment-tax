package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sathitsak/assessment-tax/pkg/handler"
)

var PERSONAL_ALLOWANCE = 60000.0

func main() {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
	port := os.Getenv("PORT")
	
	e := echo.New()
	h := handler.CreateHandler()
	e.POST("/tax/calculations", h.CalTaxHandler)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
