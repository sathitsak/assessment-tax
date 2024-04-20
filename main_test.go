package main

import "testing"

func TestCalLadder(t *testing.T){
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
		want := test.want
		got := calNetIncomeTax(test.income)
		if want != got {
			t.Errorf("%s Expect \n%v\n, got \n%v",test.name, want, got)
		}
	}
}

func TestCalTax(t *testing.T){
	t.Parallel()
	tests := []struct{
		income float64
		want float64

	}{
		{500000.0,29000.0},
		{400000.0,19000.0},
	}
	for _,test := range tests{
		want := test.want
		got := calTax(test.income)
		if want != got {
			t.Errorf(" Expect \n%v\n, got \n%v", want, got)
		}
	}
}