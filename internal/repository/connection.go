package repository

import (
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/entity"
	"github.com/jmoiron/sqlx"
)

type ConnectionRepository struct {
	db *sqlx.DB
}

func (c ConnectionRepository) GetLastConnectionPortCount() (*entity.ConnectionPortCount, error) {
	//TODO implement me
	panic("implement me")
}

func (c ConnectionRepository) CreateConnection(connection *entity.Connection) error {
	//TODO implement me
	panic("implement me")
}

func (c ConnectionRepository) GetConnectionsByUserId(id int64) (*[]entity.Connection, error) {
	//TODO implement me
	panic("implement me")
}

func NewConnectionRepository(db *sqlx.DB) *ConnectionRepository {
	return &ConnectionRepository{db}
}
