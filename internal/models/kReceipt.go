package models

import (
	"database/sql"
	"fmt"
)

type KReceiptInterface interface {
	Create(amount float64) error
	Read() (float64, error)
}

type KReceiptModel struct {
	DB *sql.DB
}

var defaultAmount = 50000.0

func (k *KReceiptModel) Create(amount float64) error {

	query := fmt.Sprintf("INSERT INTO k_receipt (amount) VALUES (%f);", amount)
	_, err := k.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil

}

func (k *KReceiptModel) Read() (float64, error) {

	var amount float64
	query := `SELECT amount FROM k_receipt  ORDER BY created_at DESC LIMIT 1;`
	row := k.DB.QueryRow(query)
	err := row.Scan(&amount)
	if err != nil {
		if err == sql.ErrNoRows {
			return defaultAmount, nil
		} else {
			return amount, err
		}
	}
	return amount, err
}
