package service

import "jwt_token/pkg"

type Authorization interface {
	CreateUser(user pkg.User) (int, error)
	AuthenticateUser(email, password string) (string, error)
}

type Service struct {
	Authorization
}

func NewService() *Service {
	auth := NewAuthService()
	return &Service{
		Authorization: auth,
	}
}
