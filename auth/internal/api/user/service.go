package user

import (
	"github.com/ivankornilov/auth/internal/service"
	desc "github.com/ivankornilov/auth/pkg/auth_v1"
)

type Implementation struct {
	desc.UnimplementedAuthV1Server
	userService service.UserService
	authService service.AuthService
}

func NewImplementation(userService service.UserService, authService service.AuthService) *Implementation {
	return &Implementation{
		userService: userService,
		authService: authService,
	}
}
