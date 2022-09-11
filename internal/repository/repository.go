package repository

import (
	"github.com/evgeniy-krivenko/particius-vpn-bot/entity"
	"github.com/jmoiron/sqlx"
)

type User interface {
	GetUserById(id int) (*entity.User, error)
}

type Repository struct {
	User
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}
