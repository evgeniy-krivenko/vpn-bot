package usecases

import (
	entity "github.com/evgeniy-krivenko/particius-vpn-bot/internal/entity"
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
}

type ServerRepository interface {
	GetAllServers() ([]entity.Server, error)
	GetServerById(id int) (*entity.Server, error)
}

type Repository interface {
	UserRepository
	TextRepository
	ConnectionRepository
	ServerRepository
}

type CryptoService interface {
	Encrypt(text []byte, key []byte) ([]byte, error)
	Decrypt(cipherText []byte, key []byte) ([]byte, error)
	GeneratePassword(passwordLen int) string
	GenerateConfig(conn *entity.Connection) (string, error)
}

type Service interface {
	CryptoService
}
