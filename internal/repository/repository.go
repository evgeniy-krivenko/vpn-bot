package repository

import (
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/usecases"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	usecases.UserRepository
	usecases.TextRepository
	usecases.ConnectionRepository
	usecases.ServerRepository
}

func New(db *sqlx.DB) usecases.Repository {
	return &Repository{
		UserRepository:       NewUserRepository(db),
		TextRepository:       NewTextRepository(db),
		ConnectionRepository: NewConnectionRepository(db),
		ServerRepository:     NewServerRepository(db),
	}
}
