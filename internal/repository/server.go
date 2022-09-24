package repository

import (
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/entity"
	"github.com/jmoiron/sqlx"
)

type ServerRepository struct {
	db *sqlx.DB
}

func (s ServerRepository) GetAllServers() (*[]entity.Server, error) {
	//TODO implement me
	panic("implement me")
}

func NewServerRepository(db *sqlx.DB) *ServerRepository {
	return &ServerRepository{db}
}
