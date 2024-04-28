// +build integration

package main
import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sathitsak/assessment-tax/internal"
	"github.com/sathitsak/assessment-tax/middleware"
	"github.com/sathitsak/assessment-tax/pkg/handler"
	"github.com/stretchr/testify/assert"
)

type IncomeData struct {
	TotalIncome float64 `json:"totalIncome"`
	WHT         float64 `json:"wht"`
	Donation    float64 `json:"donation"`
	KReceipt    float64 `json:"k-receipt"`
}

type Tax struct {
	TotalIncome float64 `json:"totalIncome" form:"totalIncome"`
	Tax         float64 `json:"tax" form:"tax"`
	TaxRefund   float64 `json:"taxRefund" form:"taxRefund"`
}

type FileUploadResponse struct {
	Taxes []Tax `json:"taxes" form:"taxes"`
}

type Decimal float64
type TaxLevel struct {
	Level string  `json:"level"`
	Tax   Decimal `json:"tax"`
}
type Response struct {
	Tax       Decimal    `json:"tax" form:"tax"`
	TaxRefund Decimal    `json:"taxRefund" form:"taxRefund"`
	TaxLevel  []TaxLevel `json:"taxLevel"`
}

type FileUploadTestCase struct {
	csvContent     string
	want           FileUploadResponse
	wantStatusCode int
}

type PersonalAllowanceResponse struct {
	PersonalDeduction float64 `json:"personalDeduction" form:"personalDeduction"`
}

type KReceiptResponse struct {
	KReceipt float64 `json:"kReceipt" form:"kReceipt"`
}

func TestAdminWrongCredential(t *testing.T) {
	var requestJSON = `{}`
	req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", strings.NewReader(requestJSON))
	auth := base64.StdEncoding.EncodeToString([]byte("adminTaxx:adminn!"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Basic "+auth)
	rec := httptest.NewRecorder()
	db, teardown := internal.SetupTestDB(t)
	defer teardown()
	h := handler.CreateHandler(db)
	e := echo.New()
	e.Use(middleware.ValidateBasicAuth("adminTax", "admin!"))
	e.POST("/admin/deductions/personal", h.PersonalAllowanceHandler)
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCalTaxHandler(t *testing.T) {
	var requestJSON = `{
		"totalIncome": 500000.0,
		"wht": 0.0,
		"allowances": [
			{
				"allowanceType": "k-receipt",
				"amount": 200000.0
			  },
			  {
				"allowanceType": "donation",
				"amount": 100000.0
			  }
		]
	  }`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	want := Response{
		Tax: 14000.0,
		TaxLevel: []TaxLevel{
			{
				Level: "0-150,000",
				Tax:   0.0,
			},
			{
				Level: "150,001-500,000",
				Tax:   14000.0,
			},
			{
				Level: "500,001-1,000,000",
				Tax:   0.0,
			},
			{
				Level: "1,000,001-2,000,000",
				Tax:   0.0,
			},
			{
				Level: "2,000,001 ขึ้นไป",
				Tax:   0.0,
			},
		},
	}

	c := e.NewContext(req, rec)
	var got Response
	db, teardown := internal.SetupTestDB(t)
	defer teardown()
	h := handler.CreateHandler(db)

	if assert.NoError(t, h.CalTaxHandler(c), json.Unmarshal(rec.Body.Bytes(), &got)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, want, got)
	}

}

func TestFileUpload(t *testing.T) {
	e := echo.New()
	tests := []FileUploadTestCase{
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
		db, teardown := internal.SetupTestDB(t)
		defer teardown()
		h := handler.CreateHandler(db)
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
	tests := []FileUploadTestCase{
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
		db, teardown := internal.SetupTestDB(t)
		defer teardown()
		h := handler.CreateHandler(db)
		c := e.NewContext(req, rec)
		h.HandleFileUpload(c)
		if assert.NoError(t, h.HandleFileUpload(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}

	}
}

func TestPersonalAllowanceHandler(t *testing.T) {
	var requestJSON = `{
		"amount": 70000.0
	  }`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", strings.NewReader(requestJSON))
	auth := base64.StdEncoding.EncodeToString([]byte("adminTax:admin!"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Basic "+auth)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	db, teardown := internal.SetupTestDB(t)
	defer teardown()
	h := handler.CreateHandler(db)
	want := PersonalAllowanceResponse{
		PersonalDeduction: 70000.0,
	}
	var got PersonalAllowanceResponse
	if assert.NoError(t, h.PersonalAllowanceHandler(c), json.Unmarshal(rec.Body.Bytes(), &got)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, got, want)
	}

}

func TestTaxRefund(t *testing.T) {
	var requestJSON = `{
		"totalIncome": 500000.0,
		"wht": 39000.0,
		"allowances": [
		  {
			"allowanceType": "donation",
			"amount": 0.0
		  }
		]
	  }`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	want := 10000.0

	c := e.NewContext(req, rec)
	var got Response
	db, teardown := internal.SetupTestDB(t)
	defer teardown()
	h := handler.CreateHandler(db)

	if assert.NoError(t, h.CalTaxHandler(c), json.Unmarshal(rec.Body.Bytes(), &got)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, want, float64(got.TaxRefund))
	}
}

func TestSetKReceipt(t *testing.T) {
	var requestJSON = `{
		"amount": 70000.0
	  }`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	db, teardown := internal.SetupTestDB(t)
	defer teardown()
	h := handler.CreateHandler(db)
	want := KReceiptResponse{
		KReceipt: 70000.0,
	}
	var got KReceiptResponse
	if assert.NoError(t, h.KReceiptHandler(c), json.Unmarshal(rec.Body.Bytes(), &got)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, got, want)
	}

}
