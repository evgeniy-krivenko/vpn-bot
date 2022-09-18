package repository

import (
	"fmt"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/entity"
	"github.com/jmoiron/sqlx"
)

type TextRepository struct {
	db *sqlx.DB
}

func (t *TextRepository) GetTermsByOrder(orderNum int) (*entity.Term, error) {
	var term entity.Term
	query := fmt.Sprintf("SELECT * FROM %s WHERE order_number=$1", termsTable)
	if err := t.db.Get(&term, query, orderNum); err != nil {
		return nil, err
	}

	return &term, nil
}

func NewTextRepository(db *sqlx.DB) *TextRepository {
	return &TextRepository{db: db}
}
