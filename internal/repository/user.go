package repository

import (
	"fmt"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	db *sqlx.DB
}

func (ur *UserRepository) SaveUser(u *entity.User) error {
	query := fmt.Sprintf("UPDATE users SET user_id=$1, username=$2, first_name=$3, last_name=$4, chat_id=$5, chat_type=$6, is_terms_confirmed=$7 WHERE user_id=%d;", u.UserId)
	_, err := ur.db.Exec(query, u.UserId, u.Username, u.FirstName, u.LastName, u.ChatId, u.ChatType, u.IsConfirmTerms)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) CreateNewUser(u *entity.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (user_id, chat_id, username, first_name, last_name, chat_type) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;", usersTable)
	row := ur.db.QueryRow(query, u.UserId, u.ChatId, u.Username, u.FirstName, u.LastName, u.ChatType)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (ur *UserRepository) GetUserById(id int) *entity.User {
	var user entity.User
	query := fmt.Sprintf("SELECT id, user_id, chat_id, username, first_name, last_name, chat_type, is_terms_confirmed FROM %s WHERE user_id=$1 OR chat_id=$2;", usersTable)
	err := ur.db.Get(&user, query, id, id)
	if err != nil {
		logrus.Errorf("[Repository] error for id {%d}: %s", id, err.Error())
		return nil
	}
	return &user
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}
