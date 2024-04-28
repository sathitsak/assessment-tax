package tax

import (
	"testing"
)

var personalAllowance = 60000.0

func TestNetIncomeTax(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		income float64
		want   float64
	}{
		{"net income: 0", 0.0, 0.0},
		{"net income: 150,000", 150000.0, 0.0},
		{"net income: 500,000", 500000.0, 35000.0},
		{"net income: 1,000,000", 1000000.0, 110000.0},
		{"net income: 2,000,000", 2000000.0, 310000.0},
		{"net income: 3,000,000.0", 3000000.0, 660000.0},
	}
	for _, test := range tests {
		tax := CreateTax(test.income, 0, 0, 0, 0)
		want := test.want
		got := tax.NetIncomeTax()
		if want != got {
			t.Errorf("%s Expect \n%v\n, got \n%v", test.name, want, got)
		}
	}
}

func TestWth(t *testing.T) {
	t.Parallel()
	tests := []struct {
		totalIncome float64
		wth         float64
		want        float64
	}{
		{500000.0, 0, 29000.0},
		{400000.0, 0, 19000.0},
		{500000.0, 25000.0, 4000.0},
		{500000.0, 39000.0, -10000.0},
	}
	for _, test := range tests {
		want := test.want
		tax := CreateTax(test.totalIncome, test.wth, personalAllowance, 0, 0)
		got := tax.PayAble()
		if want != got {
			t.Errorf(" Expect \n%v\n, got \n%v", want, got)
		}
	}
}
func TestDonation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		totalIncome float64
		donation    float64
		want        float64
	}{
		{500000.0, 200000.0, 19000.0},
		{500000.0, 100000.0, 19000.0},
		{500000.0, 50000.0, 24000.0},
		{500000.0, 0.0, 29000.0},
	}
	for _, test := range tests {
		want := test.want
		tax := CreateTax(test.totalIncome, 0.0, personalAllowance, test.donation, 0)
		got := tax.PayAble()
		if want != got {
			t.Errorf(" Expect \n%v\n, got \n%v", want, got)
		}
	}
}
func TestKReceipt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		totalIncome float64
		kReceipt    float64
		want        float64
	}{
		{500000.0, 200000.0, 24000.0},
		{500000.0, 100000.0, 24000.0},
		{500000.0, 50000.0, 24000.0},
		{500000.0, 0.0, 29000.0},
	}
	for _, test := range tests {
		want := test.want
		tax := CreateTax(test.totalIncome, 0.0, personalAllowance, 0, test.kReceipt)
		got := tax.PayAble()
		if want != got {
			t.Errorf(" Expect \n%v\n, got \n%v", want, got)
		}
	}
}

func TestTaxLevel(t *testing.T) {
	taxLevelName := [5]string{
		"0-150,000",
		"150,001-500,000",
		"500,001-1,000,000",
		"1,000,001-2,000,000",
		"2,000,001 ขึ้นไป",
	}
	tests := []struct {
		name      string
		netIncome float64
		want      [5]float64
	}{
		{name: "0", netIncome: 0, want: [5]float64{0, 0, 0, 0, 0}},
		{name: "150000", netIncome: 150000, want: [5]float64{0, 0, 0, 0, 0}},
		{name: "500000", netIncome: 500000, want: [5]float64{0, 35000, 0, 0, 0}},
		{name: "1000000", netIncome: 1000000, want: [5]float64{0, 35000, 75000, 0, 0}},
		{name: "2000000", netIncome: 2000000, want: [5]float64{0, 35000, 75000, 200000, 0}},
		{name: "3000000", netIncome: 3000000, want: [5]float64{0, 35000, 75000, 200000, 350000}},
	}
	for _, test := range tests {
		want := test.want
		tax := CreateTax(test.netIncome, 0.0, 0, 0, 0)
		got := tax.TaxLevel()
		for i, level := range tax.TaxLevel() {
			if want[i] != level.Tax || taxLevelName[i] != level.Level {
				t.Errorf(" Expect \n%v\n, got \n%v", want, got)
			}
		}
	}
}
