package usecases

import (
	"errors"
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

func (c *ConnectionUseCase) GetConnections(usr *entity.User) (*ResponseWithKeys, error) {
	conns, err := c.repo.GetConnectionsByUserId(usr.UserId)
	resp := ResponseWithKeys{}
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			svrs, err := c.repo.GetAllServers()
			if err != nil {
				panic(err)
			}

			for _, s := range *svrs {
				resp.AddRow(
					resp.AddButton(s.Location, fmt.Sprintf("create:%d", s.Id)),
				)
			}
			resp.Msg = NoConnectionsText

			return &resp, nil
		}
		return nil, errors.New("")
	}
	return nil, nil
}
