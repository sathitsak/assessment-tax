package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sathitsak/assessment-tax/pkg/tax"
)

var PERSONAL_ALLOWANCE = 60000.0

func main() {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
	port := os.Getenv("PORT")
	
	e := echo.New()
	
	e.POST("/tax/calculations", calTaxHandler)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
type Allowance struct{
	AllowanceType string `json:"allowanceType"`
	Amount float64 `json:"amount"`
}
type Request struct{
	TotalIncome float64 `json:"totalIncome"`
	Wht float64 `json:"wht"`
	Allowances []Allowance `json:"allowances"`
  }


func (req *Request) Donation()float64{
	donation :=0.0
	for _,v := range req.Allowances {
		if v.AllowanceType == "donation" {
			donation+= v.Amount
		}
	}
	return donation
}

type Response struct{
	Tax float64 `json:"tax" form:"tax"`
	TaxRefund float64 `json:"taxRefund" form:"taxRefund"`
}

func calTaxHandler(c echo.Context) error {
	var req Request
	err := c.Bind(&req); if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	tax := tax.CreateTax(req.TotalIncome,req.Wht,PERSONAL_ALLOWANCE,req.Donation())
	if tax.PayAble() >=0 {
		return c.JSON(http.StatusOK,&Response{Tax: tax.PayAble(), TaxRefund: 0.0})
	}else{
		return c.JSON(http.StatusOK,&Response{Tax: 0.0, TaxRefund: -tax.PayAble()})
	}
	
}