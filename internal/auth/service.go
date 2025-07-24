package auth

import "errors"

type Service interface {
	Authenticate(username, password string) (string, error)
}

type SimpleAuthService struct{}

func (s *SimpleAuthService) Authenticate(username, password string) (string, error) {
	if username == "staff1" && password == "password123" {
		return "staff1-token", nil
	}
	return "", errors.New("invalid credentials")
}
