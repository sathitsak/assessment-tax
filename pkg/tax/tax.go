package tax

import "math"

var MAX_DONATION = 100000.0
var ladders = []TaxLadder{
	{Level: "0-150,000", Rate: 0.0, Max: 150000.0, Min: 0.0},
	{Level: "150,001-500,000", Rate: 0.1, Max: 500000.0, Min: 150000.0},
	{Level: "500,001-1,000,000", Rate: 0.15, Max: 1000000.0, Min: 500000.0},
	{Level: "1,000,001-2,000,000", Rate: 0.2, Max: 2000000.0, Min: 1000000.0},
	{Level: "2,000,001 ขึ้นไป", Rate: 0.35, Max: math.Inf(1), Min: 2000000.0},
}
type TaxLadder struct {
	Level string  `json:"level"`
	Rate  float64 `json:"rate"`
	Max   float64 `json:"max"`
	Min   float64 `json:"min"`
}
type Tax struct {
	totalIncome float64
	wth float64
	personalAllowance float64
	donation float64
}
type Allowance struct{
	AllowanceType string `json:"donation"`
	Amount float64 `json:"amount"`
}
func CreateTax(totalIncome float64,wth float64,personalAllowance float64,donation float64)*Tax{
	return &Tax{totalIncome: totalIncome,wth: wth,personalAllowance:personalAllowance,donation: donation}
}
func (tax *Tax)NetIncome() float64{
	return tax.totalIncome - tax.personalAllowance -math.Min(tax.donation,MAX_DONATION)
}

func (tax *Tax) NetIncomeTax() float64 {
	netIncome:= tax.NetIncome()
	res := 0.0	
	for _, ladder := range ladders {
		if netIncome >= ladder.Max {
			res += (ladder.Max - ladder.Min) * ladder.Rate
		} else {
			res += (netIncome - ladder.Min) * ladder.Rate
			return res
		}

	}
	return res
}

func (tax *Tax) PayAble() float64 {
	return tax.NetIncomeTax()-tax.wth
}

type TaxLevel struct{
	Level string `json:"level"`
	Tax float64 `jsong:"tax"`
}

func (tax *Tax) TaxLevel() []TaxLevel{
	res := []TaxLevel{}
	netIncome:= tax.NetIncome()

	for _,ladder := range ladders{
		if netIncome >= ladder.Max {
			res = append(res, TaxLevel{Level: ladder.Level,Tax: (ladder.Max - ladder.Min) * ladder.Rate})
		} else {
			if netIncome > ladder.Min{
				res = append(res, TaxLevel{Level: ladder.Level,Tax: (netIncome - ladder.Min) * ladder.Rate})
			}else{
				res = append(res, TaxLevel{Level: ladder.Level,Tax: 0})
			}
		}
	}
	return res
}