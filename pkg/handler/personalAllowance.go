package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type PersonalAllowance struct {
	Amount float64 `json:"amount" form:"amount"`
}
type PersonalAllowanceResponse struct {
	PersonalDeduction float64 `json:"personalDeduction" form:"personalDeduction"`
}

func (h *Handler) PersonalAllowanceHandler(c echo.Context) error {
	var pa PersonalAllowance
	if err := c.Bind(&pa); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	if pa.Amount > 100000.0 {
		return c.String(http.StatusBadRequest, "The amount provided exceeds the maximum allowed limit.")
	}
	if pa.Amount < 10000.0 {
		return c.String(http.StatusBadRequest, "The amount provided is below the minimum allowed limit.")
	}
	err := h.personalAllowance.Create(pa.Amount)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error please contact admin or try again later")
	}
	return c.JSON(http.StatusOK, PersonalAllowanceResponse{pa.Amount})
}

func (h *Handler) PersonalAllowanceValue(c echo.Context) (float64, error) {
	v, err := h.personalAllowance.Read()
	if err != nil {
		return 0, c.String(http.StatusInternalServerError, "can't connect to database")
	}
	return v, nil
}
