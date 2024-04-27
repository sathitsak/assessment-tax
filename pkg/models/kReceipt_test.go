package models

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/sathitsak/assessment-tax/internal"

	"github.com/stretchr/testify/assert"
)

func TestSetKReceipt(t *testing.T) {
	db, teardown := internal.SetupTestDB(t)
	want := 75000.0
	pa := KReceiptModel{db}
	pa.Create(want)
	got, err := pa.Read()
	if assert.NoError(t, err) {
		assert.Equal(t, got, want)
	}
	teardown()
}





