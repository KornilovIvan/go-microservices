package user

import (
	"github.com/ivankornilov/auth/internal/service"
	desc "github.com/ivankornilov/auth/pkg/auth_v1"
)

type Implementation struct {
	desc.UnimplementedAuthV1Server
	userService service.UserService
}

func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
