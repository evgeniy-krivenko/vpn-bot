package repository

import (
	"fmt"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/entity"
	"github.com/jmoiron/sqlx"
)

type ServerRepository struct {
	db *sqlx.DB
}

func (s *ServerRepository) GetAllServers() ([]entity.Server, error) {
	var servers []entity.Server
	query := fmt.Sprintf(getAllServersSQL, serversTable)
	rows, err := s.db.Queryx(query)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return servers, nil
		}
		return nil, err
	}

	for rows.Next() {
		var srv entity.Server
		if err := rows.StructScan(&srv); err != nil {
			// TODO logging
		}
		servers = append(servers, srv)
	}

	return servers, nil
}

func (s *ServerRepository) GetServerById(id int) (*entity.Server, error) {
	var server entity.Server
	query := fmt.Sprintf(getServerSQL, serversTable)
	if err := s.db.Get(&server, query, id); err != nil {
		return nil, err
	}

	return &server, nil
}

func NewServerRepository(db *sqlx.DB) *ServerRepository {
	return &ServerRepository{db}
}
