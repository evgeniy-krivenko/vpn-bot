package usecases

import (
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/entity"
	"github.com/evgeniy-krivenko/vpn-api/gen/conn_service"
)

type UserRepository interface {
	GetUserById(id int) *entity.User
	CreateNewUser(user *entity.User) (int, error)
	SaveUser(user *entity.User) error
}

type TextRepository interface {
	GetTermsByOrder(orderNum int) (*entity.Term, error)
}

type ConnectionRepository interface {
	GetLastConnectionPortCount() (*entity.ConnectionPortCount, error)
	CreateConnection(connection *entity.Connection) (int, error)
	GetConnectionsByUserId(id int64) ([]entity.Connection, error)
	GetConnectionById(id int) (*entity.Connection, error)
	GetConnectionByServerId(id int) (*entity.Connection, error)
	SaveConnection(conn *entity.Connection) error
}

type ServerRepository interface {
	GetAllServers() ([]entity.Server, error)
	GetServerById(id int) (*entity.Server, error)
}

type Repository interface {
	UserRepository
}

type GrpcService interface {
	conn_service.ConnectionServiceClient
}
