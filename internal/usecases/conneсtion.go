package usecases

import (
	"context"
	"fmt"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/entity"
)

type ConnectionUseCase struct {
	repo Repository
	srv  Service
}

func NewConnectionUseCase(r Repository, s Service) *ConnectionUseCase {
	return &ConnectionUseCase{repo: r, srv: s}
}

func (c *ConnectionUseCase) GetConnections(ctx context.Context, usr *entity.User) (*ResponseWithKeys, error) {
	conns, err := c.repo.GetConnectionsByUserId(usr.UserId)
	if err != nil {
		return nil, err
	}

	if len(conns) == 0 {
		return c.getResponseWithoutConnections()
	}

	return c.getResponseWithConnections(conns), nil
}

func (c *ConnectionUseCase) CreateConnection(ctx context.Context, usr *entity.User) (*ResponseWithKeys, error) {
	//TODO
}

func (c *ConnectionUseCase) getResponseWithoutConnections() (*ResponseWithKeys, error) {
	resp := ResponseWithKeys{}

	svrs, err := c.repo.GetAllServers()
	if err != nil {
		return nil, err
	}

	for _, s := range svrs {
		resp.AddRow(
			resp.AddButton(s.Location, fmt.Sprintf("create:%d", s.Id)),
		)
	}
	resp.Msg = NoConnectionsText

	return &resp, nil
}

func (c *ConnectionUseCase) getResponseWithConnections(conns []entity.Connection) *ResponseWithKeys {
	resp := ResponseWithKeys{}

	msg := ConnectionsText
	for i, conn := range conns {
		msg += fmt.Sprintf("%d.: Location: %s, status: %s\\n", i+1, conn.Location, ConnectionStatusEmoji[false])
	}

	resp.Msg = msg
	return &resp
}
