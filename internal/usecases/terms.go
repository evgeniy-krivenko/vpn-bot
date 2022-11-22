package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/evgeniy-krivenko/vpn-api/gen/conn_service"
)

type TermsUseCase struct {
	Repository
	grpc GrpcService
}

func NewTermsUseCase(r Repository, g GrpcService) *TermsUseCase {
	return &TermsUseCase{
		Repository: r,
		grpc:       g,
	}
}

func (c *TermsUseCase) TermsConfirmed(ctx context.Context, userId int) (*ResponseWithKeys, error) {
	user := c.Repository.GetUserById(userId)
	if user == nil {
		return nil, errors.New("user is not found")
	}

	user.IsConfirmTerms = true
	err := c.SaveUser(user)
	if err != nil {
		return nil, err
	}

	resp := ResponseWithKeys{}

	r, err := c.grpc.GetServers(ctx, &conn_service.GetServersReq{})
	if err != nil {
		return nil, err
	}

	for _, s := range r.GetServers() {
		resp.AddRow(
			resp.AddButton(s.Location, fmt.Sprintf("create:%d", s.Id)),
		)
	}
	resp.Msg = FirstConnectionText

	return &resp, nil
}
