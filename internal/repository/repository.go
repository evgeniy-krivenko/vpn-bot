package repository

import "github.com/jmoiron/sqlx"

type User interface {
	GetUserById()
}

type Repository struct {
	User
}

func New(db *sqlx.DB) *Repository {
	return &Repository{}
}
