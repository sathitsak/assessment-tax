package handler

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sathitsak/assessment-tax/pkg/tax"
)

type IncomeData struct {
	TotalIncome float64 `json:"totalIncome"`
	WHT         float64 `json:"wht"`
	Donation    float64 `json:"donation"`
	KReceipt float64 `json:"k-receipt"`
}

type Tax struct {
	TotalIncome float64 `json:"totalIncome" form:"totalIncome"`
	Tax         float64 `json:"tax" form:"tax"`
	TaxRefund   float64 `json:"taxRefund" form:"taxRefund"`
}

type FileUploadResponse struct {
	Taxes []Tax `json:"taxes" form:"taxes"`
}

func (h *handler) HandleFileUpload(c echo.Context) error {
	file, err := c.FormFile("taxFile")
	if err != nil {
		return  c.String(http.StatusBadRequest, "Failed to get file from form data")
	}

	src, err := file.Open()
	if err != nil {
		return  c.String(http.StatusBadRequest, "Failed to open file")
	}
	defer src.Close()

	csvReader := csv.NewReader(src)

	headers, err := csvReader.Read()
	if err != nil {
		return  c.String(http.StatusBadRequest, "Failed to read headers from CSV")
	}
	if len(headers) != 3 && len(headers) != 4{
		return  c.String(http.StatusBadRequest, "CSV does not contain the required headers")
	}
	if len(headers) == 3  && (headers[0] != "totalIncome" || headers[1] != "wht" || headers[2] != "donation") {
		return  c.String(http.StatusBadRequest, "CSV does not contain the required headers")
	}
	if len(headers) == 4  && (headers[0] != "totalIncome" || headers[1] != "wht" || headers[2] != "donation" || headers[3] != "k-receipt") {
		return  c.String(http.StatusBadRequest, "CSV does not contain the required headers")
	}

	var data []IncomeData

	// Read each row and convert it into an IncomeData object
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return c.String(http.StatusBadRequest, "Failed to read data row")
		}

		var rowData IncomeData
		rowData.TotalIncome, err = strconv.ParseFloat(row[0], 64)
		if err != nil {
			return  c.String(http.StatusBadRequest, fmt.Sprintf("Invalid totalIncome value: %v", row[0]))
		}
		rowData.WHT, err = strconv.ParseFloat(row[1], 64)
		if err != nil {
			return  c.String(http.StatusBadRequest, fmt.Sprintf("Invalid WHT value: %v", row[1]))
		}
		rowData.Donation, err = strconv.ParseFloat(row[2], 64)
		if err != nil {
			return  c.String(http.StatusBadRequest, fmt.Sprintf("Invalid Donation value: %v", row[2]))
		}
		rowData.KReceipt = 0.0
		if len(headers) == 4 {
			rowData.KReceipt, err = strconv.ParseFloat(row[3], 64)
		if err != nil {
			return  c.String(http.StatusBadRequest, fmt.Sprintf("Invalid KReceipt value: %v", row[3]))
		}
		}

		data = append(data, rowData)

	}
	pa, err := h.personalAllowance.Read()
	if err != nil {
		return c.String(http.StatusBadRequest, "Internal server error please contact admin or try again later")
	}
	res := []Tax{}
	for _, row := range data {
		tax := tax.CreateTax(row.TotalIncome, row.WHT, pa, row.Donation,row.KReceipt)
		if tax.PayAble() >= 0 {
			res = append(res, Tax{TotalIncome: row.TotalIncome, Tax: tax.PayAble(), TaxRefund: 0})
		} else {
			res = append(res, Tax{TotalIncome: row.TotalIncome, Tax: 0, TaxRefund: tax.PayAble()})
		}
	}

	// Optionally, you can print or process `data` further here
	fmt.Printf("Received data: %+v\n", data)

	return c.JSON(http.StatusOK, FileUploadResponse{Taxes: res})
}
