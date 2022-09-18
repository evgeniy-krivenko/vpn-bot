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

type Repository interface {
	UserRepository
	TextRepository
}

type KeyboardService interface {
	GetInlineKeyboardRow()
	GetInlineKeyboardButtonData()
}
