package handler

import (
	"database/sql"
	"fmt"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sathitsak/assessment-tax/pkg/models"
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

type handler struct{
	personalAllowance models.PersonalAllowanceInterface
}

func (req *Request) Donation() float64 {
	donation := 0.0
	for _, v := range req.Allowances {
		if v.AllowanceType == "donation" {
			donation += v.Amount
		}
	}
	return donation
}

func (req *Request) KReceipt() float64 {
	kReceipt := 0.0
	for _, v := range req.Allowances {
		if v.AllowanceType == "k-receipt" {
			kReceipt += v.Amount
		}
	}
	
	return kReceipt
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
func CreateHandler(db *sql.DB)handler{
	return handler{
		personalAllowance: &models.PersonalAllowanceModel{DB: db},
	}
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
	
	
	pa,err := h.personalAllowance.Read()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error please contact admin or try again later")
	}
	tax := tax.CreateTax(req.TotalIncome, req.Wht, pa, req.Donation(),req.KReceipt())
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


type PersonalAllowance struct{
	Amount float64 `json:"amount" form:"amount"`
}
type PersonalAllowanceResponse struct{
	PersonalDeduction float64 `json:"personalDeduction" form:"personalDeduction"`
}

func (h *handler) PersonalAllowanceHandler (c echo.Context) error {
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
	err := h.personalAllowance.Create(pa.Amount)
	if  err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error please contact admin or try again later")
	}
	return c.JSON(http.StatusOK,PersonalAllowanceResponse{pa.Amount})
}

