package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"fmt"
	"math"
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
	
	// e.POST("/tax/calculations", calTax)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
type Allowance struct{
	AllowanceType string `json:"donation"`
	Amount float64 `json:"amount"`
}
type Tax struct{
	TotalIncome float64 `json:"totalIncome"`
	Wht float64 `json:"wht"`
	Allowances []Allowance `json:"allowances"`
  }


func calNetIncomeTax(income float64)float64{
	res:=0.0
	taxLadders := []TaxLadder{
		{Level: "0-150,000" ,Rate: 0.0, Max:150000.0, Min:0.0},
		{Level: "150,001 - 500,000",Rate:0.1,Max: 500000.0,Min:150000.0},
		{Level: "500,001 - 1,000,000",Rate: 0.15,Max: 1000000.0,Min:500000.0},
		{Level: "1,000,001 - 2,000,000",Rate: 0.2, Max: 2000000.0,Min:1000000.0},
		{Level: "2,000,001 ขึ้นไป",Rate: 0.35,Max: math.Inf(1),Min:2000000.0},

	}
	for _,ladder := range taxLadders{
		if income >= ladder.Max {
			res+= (ladder.Max - ladder.Min)*ladder.Rate
		} else {
			res += (income-ladder.Min)*ladder.Rate
			return res
		}
		
	}
	return res
}
func calTax(totalIncome float64)float64{
	return calNetIncomeTax(totalIncome-60000.0)
}


type TaxLadder struct{
	Level string `json:"level"`
	Rate float64 `json:"rate"`
	Max float64 `json:"max"`
	Min float64 `json:"min"`
}


// func calTax(c echo.Context) error {
// 	// User ID from path `users/:id`
// 	id := c.Param("id")
//   return c.String(http.StatusOK, id)
// }