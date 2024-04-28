package handler

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	csvContent     string
	want           FileUploadResponse
	wantStatusCode int
}

func TestFileUpload(t *testing.T) {
	e := echo.New()
	tests := []TestCase{
		{csvContent: "totalIncome,wht,donation\n500000,0,0\n600000,40000,20000\n750000,50000,15000",
			want: FileUploadResponse{
				Taxes: []Tax{
					{TotalIncome: 500000, Tax: 29000, TaxRefund: 0.0},
					{TotalIncome: 600000, Tax: 0, TaxRefund: -2000.0},
					{TotalIncome: 750000, Tax: 11250, TaxRefund: 0},
				},
			},
			wantStatusCode: http.StatusOK},
		{csvContent: "totalIncome,wht,donation,k-receipt\n500000,0,0,0\n500000,0,100000.0,200000.0",
			want: FileUploadResponse{Taxes: []Tax{
				{TotalIncome: 500000, Tax: 29000, TaxRefund: 0.0},
				{TotalIncome: 500000, Tax: 14000, TaxRefund: 0},
				// {TotalIncome: 750000, Tax: 0, TaxRefund: 38750.0},
			}}, wantStatusCode: http.StatusOK,
		},
	}
	for _, test := range tests {

		// Create a buffer to write our multipart form data
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		// Simulate a file with CSV content
		part, err := writer.CreateFormFile("taxFile", "test.csv")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := part.Write([]byte(test.csvContent)); err != nil {
			t.Fatal(err)
		}
		writer.Close()

		// Create a request to your endpoint
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Record the response
		rec := httptest.NewRecorder()
		h := CreateTestHandler()
		c := e.NewContext(req, rec)
		var got FileUploadResponse

		if assert.NoError(t, h.HandleFileUpload(c), json.Unmarshal(rec.Body.Bytes(), &got)) {
			assert.Equal(t, test.wantStatusCode, rec.Code)
			assert.Equal(t, test.want, got)
		}
	}
}
func TestFileUploadInvalidFile(t *testing.T) {
	e := echo.New()
	tests := []TestCase{
		{csvContent: "totalIncome,wht,donation\n500000,,0\n600000,,20000\n750000,50000,15000",
			want: FileUploadResponse{
				Taxes: []Tax{
					{TotalIncome: 500000, Tax: 29000, TaxRefund: 0.0},
					{TotalIncome: 600000, Tax: 0, TaxRefund: -2000.0},
					{TotalIncome: 750000, Tax: 11250, TaxRefund: 0},
				},
			},
			wantStatusCode: http.StatusOK},
		{csvContent: "totalIncome,wht,,k-receipt\n500000,0,0,0\n500000,0,100000.0,200000.0",
			want: FileUploadResponse{Taxes: []Tax{
				{TotalIncome: 500000, Tax: 29000, TaxRefund: 0.0},
				{TotalIncome: 500000, Tax: 14000, TaxRefund: 0},
				// {TotalIncome: 750000, Tax: 0, TaxRefund: 38750.0},
			}}, wantStatusCode: http.StatusOK,
		}, {csvContent: "totalIncome,wht,donation,k-receipt\n500000,0,0\n600000,,50000\n750000,50000,100000",
			want: FileUploadResponse{}, wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, test := range tests {

		// Create a buffer to write our multipart form data
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		// Simulate a file with CSV content
		part, err := writer.CreateFormFile("taxFile", "test.csv")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := part.Write([]byte(test.csvContent)); err != nil {
			t.Fatal(err)
		}
		writer.Close()

		// Create a request to your endpoint
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Record the response
		rec := httptest.NewRecorder()
		h := CreateTestHandler()
		c := e.NewContext(req, rec)
		h.HandleFileUpload(c)
		if assert.NoError(t, h.HandleFileUpload(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}

	}
}
