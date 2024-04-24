package repositories

import (
	"database/sql"
	"fmt"
)

func CreatePersonalAllowance(db *sql.DB, amount float64) (int64, error) {
	
	query := fmt.Sprintf("INSERT INTO personal_allowance (amount) VALUES (%f);", amount)
	res, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()

}

func ReadPersonalAllowance(db *sql.DB) (float64, error){
	
	var amount float64
	query := `SELECT amount FROM personal_allowance  ORDER BY created_at DESC LIMIT 1;`
	row := db.QueryRow(query)
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