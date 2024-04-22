package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sathitsak/assessment-tax/pkg/tax"
)
var PERSONAL_ALLOWANCE = 60000.0

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}
type Request struct {
	TotalIncome float64     `json:"totalIncome"`
	Wht         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type handler struct{}

func (req *Request) Donation() float64 {
	donation := 0.0
	for _, v := range req.Allowances {
		if v.AllowanceType == "donation" {
			donation += v.Amount
		}
	}
	return donation
}

type Response struct {
	Tax       Decimal    `json:"tax" form:"tax"`
	TaxRefund Decimal    `json:"taxRefund" form:"taxRefund"`
	TaxLevel  []TaxLevel `json:"taxLevel"`
}
type Decimal float64
type TaxLevel struct {
	Level string  `json:"level"`
	Tax   Decimal `json:"tax"`
}
func CreateHandler()handler{
	return handler{}
}
func (d Decimal) MarshalJSON() ([]byte, error) {
	// Always format with one decimal place
	return []byte(fmt.Sprintf("%.1f", d)), nil
}
func (h *handler)CalTaxHandler(c echo.Context) error {
	var req Request
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	tax := tax.CreateTax(req.TotalIncome, req.Wht, PERSONAL_ALLOWANCE, req.Donation())
	taxLevel := []TaxLevel{}
	for _, v := range tax.TaxLevel() {
		taxLevel = append(taxLevel, TaxLevel{Level: v.Level, Tax: Decimal(v.Tax)})
	}
	if tax.PayAble() >= 0 {
		return c.JSON(http.StatusOK, &Response{Tax: Decimal(tax.PayAble()), TaxRefund: 0.0, TaxLevel: taxLevel})
	} else {
		return c.JSON(http.StatusOK, &Response{Tax: 0.0, TaxRefund: Decimal(-tax.PayAble()), TaxLevel: taxLevel})
	}

}