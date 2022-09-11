package service

import (
	"context"
	"fmt"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/repository"
)

type Service struct {
}

func NewService(r *repository.Repository) *Service {
	return &Service{}
}

// Start получить данные о юзере, которые мы сохраним в базе
func (s *Service) Start(ctx context.Context, dto StartDto) string {
	fmt.Printf("Chat ID: %d\n", dto.ChatId)
	fmt.Printf("FirstNane: %s\n", dto.FirstName)
	fmt.Printf("UserId: %d\n", dto.UserId)
	return "Спасибо за регистрацию!"
}
