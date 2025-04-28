package service

import (
	"errors"
	"jwt_token/pkg"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users        []pkg.User
	tokenManager TokenManager
}

func NewAuthService() *AuthService {
	return &AuthService{
		users:        []pkg.User{},
		tokenManager: NewJWTManager(signingKey),
	}
}

const customSalt = "MySuperSecretSalt123"

func hashPassword(password string) (string, error) {
	salted := password + customSalt
	bytes, err := bcrypt.GenerateFromPassword([]byte(salted), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	salted := password + customSalt
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(salted))
	return err == nil
}

func (s *AuthService) CreateUser(user pkg.User) (int, error) {
	for _, u := range s.users {
		if u.Email == user.Email {
			return 0, errors.New("user already exists")
		}
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return 0, err
	}

	user.Id = len(s.users) + 1
	user.Password = hashedPassword
	s.users = append(s.users, user)
	return user.Id, nil
}

func (s *AuthService) AuthenticateUser(email, password string) (string, error) {
	for _, u := range s.users {
		if u.Email == email && checkPasswordHash(password, u.Password) {
			return s.tokenManager.GenerateToken(u.Id)
		}
	}
	return "", errors.New("invalid email or password")
}
