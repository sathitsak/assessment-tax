package handler

import (
	"database/sql"
	"fmt"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sathitsak/assessment-tax/internal/models"
	"github.com/sathitsak/assessment-tax/internal/tax"
)

var PERSONAL_ALLOWANCE = 60000.0

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type Request struct {
	TotalIncome *float64     `json:"totalIncome"`
	Wht         *float64     `json:"wht"`
	Allowances  *[]Allowance `json:"allowances"`
}

type Handler struct {
	personalAllowance models.PersonalAllowanceInterface
	kReceipt          models.KReceiptInterface
}

func Donation(req *Request) float64 {
	donation := 0.0
	for _, v := range *req.Allowances {
		if v.AllowanceType == "donation" {
			donation += v.Amount
		}
	}
	return donation
}

func KReceipt(req *Request) float64 {
	kReceipt := 0.0
	for _, v := range *req.Allowances {
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

func CreateHandler(db *sql.DB) Handler {
	return Handler{
		personalAllowance: &models.PersonalAllowanceModel{DB: db},
		kReceipt:          &models.KReceiptModel{DB: db},
	}
}
func (d Decimal) MarshalJSON() ([]byte, error) {
	// Always format with one decimal place
	return []byte(fmt.Sprintf("%.1f", d)), nil
}
func (h *Handler) CalTaxHandler(c echo.Context) error {
	req := new(Request)

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	// Check if totalIncome and wht are provided
	if req.TotalIncome == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "totalIncome is required"})
	}
	if req.Wht == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "wht is required"})
	}

	// Check if allowances is provided and valid
	if req.Allowances == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "allowances are required"})
	}
	for _, allowance := range *req.Allowances {
		if allowance.AllowanceType != "k-receipt" && allowance.AllowanceType != "donation" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "unknown allowanceType"})
		}
	}

	pa, err := h.personalAllowance.Read()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error please contact admin or try again later")
	}
	tax := tax.CreateTax(*req.TotalIncome, *req.Wht, pa, Donation(req), KReceipt(req))
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
