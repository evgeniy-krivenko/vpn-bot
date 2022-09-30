package repository

import (
	"fmt"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ConnectionRepository struct {
	db *sqlx.DB
}

func (c *ConnectionRepository) GetLastConnectionPortCount() (*entity.ConnectionPortCount, error) {
	var cp entity.ConnectionPortCount
	query := fmt.Sprintf(getPortWithCountSQL, connectionsTable)
	if err := c.db.Get(&cp, query); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	return &cp, nil
}

func (c *ConnectionRepository) GetConnectionById(id int) (*entity.Connection, error) {
	var conn entity.Connection
	query := fmt.Sprintf(
		`SELECT c.id, s.location, c.port, c.user_id, c.encrypted_secret, s.ip_address,
						c.server_id, c.is_active, c.last_activate
		FROM %s c LEFT JOIN %s s on s.id=c.server_id WHERE c.id=$1`,
		connectionsTable,
		serversTable)
	if err := c.db.Get(&conn, query, id); err != nil {
		return nil, err
	}

	return &conn, nil
}

func (c *ConnectionRepository) CreateConnection(conn *entity.Connection) (int, error) {
	var id int
	query := fmt.Sprintf(createConnectionSQL, connectionsTable)
	row := c.db.QueryRowx(query, conn.Port, conn.EncryptedSecret, conn.UserId, conn.ServerId)
	if err := row.Scan(&id); err != nil {
		logrus.Errorf("error when creating conn: %s, conn: %v", err.Error(), conn)
		return 0, err
	}
	return id, nil
}

func (c *ConnectionRepository) GetConnectionsByUserId(id int64) ([]entity.Connection, error) {
	var connections []entity.Connection

	query := fmt.Sprintf(GetConnectsByUserId, connectionsTable)
	rows, err := c.db.Queryx(query, id)
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			logrus.Errorf("error when closing conn: %s", err.Error())
		}
	}(rows)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return connections, nil
		}
		return nil, err
	}

	for rows.Next() {
		var conn entity.Connection
		if err := rows.StructScan(&conn); err != nil {
			logrus.Errorf("error scan row to struct: %s", err.Error())
			continue
		}
		connections = append(connections, conn)
	}
	return connections, nil
}

func NewConnectionRepository(db *sqlx.DB) *ConnectionRepository {
	return &ConnectionRepository{db}
}
