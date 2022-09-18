package repository

import (
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/usecases"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	usecases.UserRepository
	usecases.TextRepository
}

func New(db *sqlx.DB) usecases.Repository {
	return &Repository{
		UserRepository: NewUserRepository(db),
		TextRepository: NewTextRepository(db),
	}
}
