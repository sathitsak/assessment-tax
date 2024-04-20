package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"fmt"
	"github.com/sathitsak/assessment-tax/pkg/tax"
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
	
	e.POST("/tax/calculations", calTaxHandler)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
type Allowance struct{
	AllowanceType string `json:"donation"`
	Amount float64 `json:"amount"`
}
type Request struct{
	TotalIncome float64 `json:"totalIncome"`
	Wht float64 `json:"wht"`
	Allowances []Allowance `json:"allowances"`
  }

type Response struct{
	Tax float64 `json:"tax" form:"tax"`
}

func calTaxHandler(c echo.Context) error {
	var req Request
	err := c.Bind(&req); if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	res := &Response{
		Tax: tax.CalTax(req.TotalIncome),
	}
	
  return c.JSON(http.StatusOK, res)
}