package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type IncomeRequest struct {
	TotalIncome *float64     `json:"totalIncome"`
	WHT         *float64     `json:"wht"`
	Allowances  *[]Allowance `json:"allowances"`
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

func ValidateRequestMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req IncomeRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
		}

		// Check if totalIncome and wht are provided
		if req.TotalIncome == nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "totalIncome is required"})
		}
		if req.WHT == nil {
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

		// Proceed with the next handler
		return next(c)
	}
}
