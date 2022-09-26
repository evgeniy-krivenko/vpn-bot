package usecases

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/entity"
	"os"
)

const MaxConnectionByPort = 10
const SecretLength = 12

// TODO move to env
var SecretKey = os.Getenv("SECRET_KEY")

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

func (c *ConnectionUseCase) CreateConnection(ctx context.Context, usr *entity.User, serverId int) ([]ResponseWithKeys, error) {
	server, err := c.repo.GetServerById(serverId)
	if err != nil {
		return nil, err
	}

	cp, err := c.repo.GetLastConnectionPortCount()
	if err != nil {
		return nil, err
	}

	if cp == nil {
		// TODO move to env default port
		cp = &entity.ConnectionPortCount{Port: 3020, Count: 0}
	}

	if cp.Count >= MaxConnectionByPort {
		cp.Port += 1
		// TODO check if port opened or use other port
	}

	plainSecret := c.srv.GeneratePassword(SecretLength)
	encryptedSecret, err := c.srv.Encrypt([]byte(plainSecret), []byte(SecretKey))
	if err != nil {
		return nil, err
	}

	conn := entity.Connection{
		Location:        server.Location,
		IpAddress:       server.IpAddress,
		UserId:          usr.Id,
		EncryptedSecret: hex.EncodeToString(encryptedSecret),
		Port:            cp.Port,
		ServerId:        server.Id,
	}

	cnf, err := c.srv.GenerateConfig(&conn)
	if err != nil {
		return nil, err
	}

	id, err := c.repo.CreateConnection(&conn)
	if err != nil {
		return nil, err
	}

	rspWithKeys := ResponseWithKeys{Msg: ConnectCreated}
	rspWithKeys.AddRow(
		rspWithKeys.AddButton(fmt.Sprintf("Активировать %s", conn.Location), fmt.Sprintf("activate:%d", id)),
	)

	var result []ResponseWithKeys
	result = append(result, ResponseWithKeys{Msg: cnf}, rspWithKeys)
	return result, nil
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
		msg += fmt.Sprintf("%d\\. Location: %s, status: %s\n", i+1, conn.Location, ConnectionStatusEmoji[false])
	}

	resp.Msg = msg
	return &resp
}
