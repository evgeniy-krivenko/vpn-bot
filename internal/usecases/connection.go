package usecases

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/entity"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/logger"
	"github.com/evgeniy-krivenko/particius-vpn-bot/pkg/e"
	"github.com/evgeniy-krivenko/particius-vpn-bot/utils"
	"github.com/evgeniy-krivenko/vpn-api/gen/conn_service"
	"time"
)

type ConnectionUseCase struct {
	repo Repository
	grpc GrpcService
	log  logger.Logger
}

func NewConnectionUseCase(r Repository, grpc GrpcService, l logger.Logger) *ConnectionUseCase {
	return &ConnectionUseCase{repo: r, grpc: grpc, log: l}
}

// GetConnections - переписано на vpn-asynq
func (c *ConnectionUseCase) GetConnections(ctx context.Context, usr *entity.User) (*ResponseWithKeys, error) {
	res, err := c.grpc.GetConnections(ctx, &conn_service.GetConnectionsReq{UserId: usr.UserId})
	if err != nil {
		return nil, err
	}

	if len(res.Connections) == 0 {
		return c.getResponseWithoutConnections(ctx)
	}

	return c.getResponseWithConnections(res.Connections), nil
}

// CreateConnection - переписано на vpn-asynq
func (c *ConnectionUseCase) CreateConnection(ctx context.Context, usr *entity.User, serverId int) ([]ResponseWithKeys, error) {
	var result []ResponseWithKeys
	// TODO получить коннект по сервер id
	conns, err := c.grpc.GetConnections(ctx, &conn_service.GetConnectionsReq{UserId: usr.UserId})
	if err != nil {
		return nil, e.Warp(fmt.Sprintf(
			"error to get connections for userId: %d",
			usr.UserId,
		), err)
	}

	var isConnectionExists bool
	for _, c := range conns.Connections {
		if c.UserId == usr.UserId && (int(c.ServerId) == serverId) {
			isConnectionExists = true
		}
	}

	if isConnectionExists {
		r := ResponseWithKeys{Msg: ConnectionExists, IsMessageDelete: true}
		return append(result, r), nil
	}

	newConn, err := c.grpc.CreateConnection(
		ctx,
		&conn_service.CreateConnectionReq{UserId: usr.UserId, ServerId: int64(serverId)},
	)
	if err != nil {
		return nil, err
	}

	cnf, err := c.grpc.GetConfig(ctx, &conn_service.GetConfigReq{
		UserId:       usr.UserId,
		ConnectionId: newConn.Id,
	})
	if err != nil {
		return nil, err
	}

	rspWithKeys := ResponseWithKeys{Msg: ConnectCreated}
	rspWithKeys.AddRow(
		rspWithKeys.AddButton(
			fmt.Sprintf("Активировать %s", newConn.Location),
			fmt.Sprintf("show-ads:%d", newConn.Id),
		),
	)

	configMessage := fmt.Sprintf("`%s`", cnf.GetConfig())

	result = append(
		result,
		ResponseWithKeys{Msg: configMessage},
		rspWithKeys,
	)
	return result, nil
}

func (c *ConnectionUseCase) OpenConnection(ctx context.Context, id int) (*ResponseWithKeys, error) {
	connResp, err := c.grpc.GetConnectionInfo(ctx, &conn_service.GetConnectionInfoReq{
		Id: int64(id),
	})
	if err != nil {
		return nil, err
	}

	lastActivationTime := sql.NullTime{}

	err = lastActivationTime.Scan(connResp.LastActivate.AsTime())
	if err != nil {
		c.log.WithContextReqId(ctx).
			Warn(fmt.Errorf(
				"error to scan time for conn id %d: %w",
				connResp.Id,
				err,
			))
	}

	conn := entity.Connection{
		Id:           int(connResp.Id),
		Location:     connResp.Location,
		Port:         uint(connResp.Port),
		UserId:       int(connResp.UserId),
		IpAddress:    connResp.IpAddress,
		ServerId:     int(connResp.ServerId),
		IsActive:     connResp.IsActive,
		LastActivate: lastActivationTime,
	}

	resp := ResponseWithKeys{}
	resp.Msg = fmt.Sprintf(
		ConnectionInfo,
		LocationFullName[conn.Location[:2]],
		ConnectionStatusText[conn.IsActive],
	)

	if conn.IsActive {
		timeLeft := conn.LastActivate.Time.Add(time.Hour * 24).Sub(time.Now())
		resp.Msg += fmt.Sprintf(
			ConnectionTimeLeft,
			utils.CutTimeString(timeLeft.String()),
		)
	} else {
		resp.AddRow(
			resp.AddButton(ActivateBtn, fmt.Sprintf("show-ads:%d", conn.Id)),
		)
	}
	resp.AddRow(
		resp.AddButton("Получить конфигурацию", fmt.Sprintf("get-config:%d", conn.Id)),
	)

	// TODO реализовать логику удаления в отдельном методе
	resp.AddRow(
		resp.AddButton("Удалить", fmt.Sprintf("delete-with-confirm:%d", conn.Id)),
	)

	return &resp, nil
}

// ActivateConnection - переписано на vpn-asynq
func (c *ConnectionUseCase) ActivateConnection(ctx context.Context, id int) (*ResponseWithKeys, error) {
	conn, err := c.grpc.GetConnectionInfo(ctx, &conn_service.GetConnectionInfoReq{Id: int64(id)})
	if err != nil {
		return nil, err
	}

	if conn.IsActive {
		return &ResponseWithKeys{
			Msg:             "Подключение уже активно",
			IsMessageDelete: false,
		}, nil
	}

	_, err = c.grpc.ActivateConnection(
		ctx,
		&conn_service.SwitchConnectionReq{Id: conn.Id},
	)
	if err != nil {
		return nil, err
	}

	return &ResponseWithKeys{
		Msg:             "Подключение активно",
		IsMessageDelete: false,
	}, nil
}

func (c *ConnectionUseCase) GetConfiguration(ctx context.Context, id int) ([]ResponseWithKeys, error) {
	var result []ResponseWithKeys
	conf, err := c.grpc.GetConfig(ctx, &conn_service.GetConfigReq{
		ConnectionId: int64(id),
		UserId:       0,
	})
	if err != nil {
		return nil, err
	}

	result = append(result,
		ResponseWithKeys{Msg: ShowConnectionConfig, IsMessageDelete: true},
		ResponseWithKeys{Msg: fmt.Sprintf("`%s`", conf.GetConfig()), IsMessageDelete: true})

	return result, nil
}

func (c *ConnectionUseCase) ShowAds(ctx context.Context, usr *entity.User, id int) (*ResponseWithKeys, error) {
	conn, err := c.grpc.GetConnectionInfo(ctx, &conn_service.GetConnectionInfoReq{Id: int64(id)})
	if err != nil {
		return nil, err
	}

	if conn.IsActive {
		return &ResponseWithKeys{
			Msg:             "Подключение уже активно",
			IsMessageDelete: true,
		}, nil
	}

	rspWithKeys := ResponseWithKeys{Msg: AdvertisingMock}
	rspWithKeys.AddRow(
		rspWithKeys.AddButton(
			ActivateDoneBtn,
			fmt.Sprintf("activate:%d", id),
		),
	)
	return &rspWithKeys, nil
}

func (c *ConnectionUseCase) getResponseWithoutConnections(ctx context.Context) (*ResponseWithKeys, error) {
	resp := ResponseWithKeys{}
	var userId int64 = 1

	r, err := c.grpc.GetServers(ctx, &conn_service.GetServersReq{UserId: &userId})
	if err != nil {
		return nil, err
	}

	for _, s := range r.GetServers() {
		resp.AddRow(
			resp.AddButton(s.Location, fmt.Sprintf("create:%d", s.Id)),
		)
	}
	resp.Msg = NoConnectionsText

	return &resp, nil
}

func (c *ConnectionUseCase) getResponseWithConnections(conns []*conn_service.Connection) *ResponseWithKeys {
	resp := ResponseWithKeys{}

	resp.Msg = ConnectionsText
	for i, conn := range conns {
		btnText := fmt.Sprintf(
			"%d. Сервер: %s, статус: %s\n",
			i+1,
			conn.Location,
			ConnectionStatusEmoji[conn.IsActive],
		)
		btnAction := fmt.Sprintf("open-connection:%d", conn.Id)
		resp.AddRow(
			resp.AddButton(btnText, btnAction),
		)
	}

	return &resp
}
