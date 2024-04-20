package tax

import (
	"testing"

)
var personalAllowance = 60000.0
func TestNetIncomeTax(t *testing.T){
	t.Parallel()
	tests := []struct{
		name string
		income float64
		want float64

	}{
		{"net income: 0",0.0,0.0},
		{"net income: 150,000",150000.0,0.0},
		{"net income: 500,000",500000.0,35000.0},
		{"net income: 1,000,000",1000000.0,110000.0},
		{"net income: 2,000,000",2000000.0,310000.0},
		{"net income: 3,000,000.0",3000000.0,660000.0},
	}
	for _,test := range tests{
		tax := CreateTax(test.income+personalAllowance,0,personalAllowance,0)
		want := test.want
		got := tax.NetIncomeTax()
		if want != got {
			t.Errorf("%s Expect \n%v\n, got \n%v",test.name, want, got)
		}
	}
}

func TestPayable(t *testing.T){
	t.Parallel()
	tests := []struct{
		totalIncome float64
		wth float64
		want float64

	}{
		{500000.0,0,29000.0},
		{400000.0,0,19000.0},
		{500000.0,25000.0,4000.0},
		{500000.0,39000.0,-10000.0},
	}
	for _,test := range tests{
		want := test.want
		tax := CreateTax(test.totalIncome,test.wth,personalAllowance,0)
		got := tax.PayAble()
		if want != got {
			t.Errorf(" Expect \n%v\n, got \n%v", want, got)
		}
	}
}
