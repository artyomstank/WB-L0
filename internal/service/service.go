package service

import "L0-wb/internal/repo"

type UserService struct {
	UserRepo repo.PostgresRepo
}

func NewService(ur repo.PostgresRepo) UserService {
	return UserService{
		UserRepo: ur,
	}
}
