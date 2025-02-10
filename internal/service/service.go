package service

import (
	"context"
	"github.com/kstsm/avito-shop/internal/repository"
)

type ServiceI interface {
	Authenticate(ctx context.Context, username, password string) (string, error)
	getApiInfo()
}

type Service struct {
	repo repository.RepositoryI
}

func NewService(repo repository.RepositoryI) ServiceI {
	return &Service{
		repo: repo,
	}
}
