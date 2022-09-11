package repository

import (
	"fmt"
	"github.com/evgeniy-krivenko/particius-vpn-bot/entity"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func (u *UserRepository) GetUserById(id int) (*entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT user_id, chat_id, username, first_name, last_name, chat_type FROM %s; WHERE user_id = %d OR chat_id = %d", usersTable, id, id)
	row := u.db.QueryRow(query)
	if err := row.Scan(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}
