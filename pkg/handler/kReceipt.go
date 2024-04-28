package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type KReceiptRequest struct {
	Amount float64 `json:"amount" form:"amount"`
}
type KReceiptResponse struct {
	KReceipt float64 `json:"kReceipt" form:"kReceipt"`
}

var minAmount = 0.0
var maxAmount = 100000.0

func (h *Handler) KReceiptHandler(c echo.Context) error {
	var k KReceiptRequest
	if err := c.Bind(&k); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	if k.Amount < minAmount {
		return c.String(http.StatusBadRequest, "The amount provided is below the minimum allowed limit.")
	}
	if k.Amount > maxAmount {
		return c.String(http.StatusBadRequest, "The amount provided exceeds the maximum allowed limit.")
	}
	err := h.kReceipt.Create(k.Amount)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error please contact admin or try again later")
	}
	return c.JSON(http.StatusOK, KReceiptResponse{k.Amount})
}
