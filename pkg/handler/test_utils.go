package handler

import "github.com/sathitsak/assessment-tax/internal/mock"

func CreateTestHandler() Handler {
	return Handler{
		personalAllowance: &mock.MockDB{Amounts: []float64{}, DefaultAmount: 60000.0},
		kReceipt:          &mock.MockDB{Amounts: []float64{}, DefaultAmount: 50000.0},
	}
}
