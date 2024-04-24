package models

import (
	"database/sql"
	"fmt"

)

type PersonalAllowance struct{
	amount float64
}
type PersonalAllowanceInterface interface{
	Create(amount float64) error
	Read()(float64, error)
}

type PersonalAllowanceModel struct {
	DB *sql.DB
}

func (pa *PersonalAllowanceModel)Create(amount float64)  error {
	
	query := fmt.Sprintf("INSERT INTO personal_allowance (amount) VALUES (%f);", amount)
	_, err := pa.DB.Exec(query)
	if err != nil {
		return  err
	}
	return nil

}

func (pa *PersonalAllowanceModel)Read() (float64, error){
	
	var amount float64
	query := `SELECT amount FROM personal_allowance  ORDER BY created_at DESC LIMIT 1;`
	row := pa.DB.QueryRow(query)
    err := row.Scan(&amount)
	if err != nil{
		if err ==  sql.ErrNoRows{
			return 60000.0,nil
		}else{
			return amount,err
		}
	}
	return amount, err
}

