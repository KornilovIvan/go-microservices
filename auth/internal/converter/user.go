package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ivankornilov/auth/internal/model"
	desc "github.com/ivankornilov/auth/pkg/auth_v1"
)

func ToUserInfoFromDesc(req *desc.CreateRequest) *model.UserInfo {
	return &model.UserInfo{
		Name:            req.GetName(),
		Email:           req.GetEmail(),
		Password:        req.GetPassword(),
		PasswordConfirm: req.GetPasswordConfirm(),
		Role:            req.GetRole().String(),
	}
}

func ToUpdateUserFromDesc(req *desc.UpdateRequest) *model.UpdateUser {
	user := &model.UpdateUser{ID: req.GetId()}
	if req.GetName() != nil {
		name := req.GetName().GetValue()
		user.Name = &name
	}
	if req.GetEmail() != nil {
		email := req.GetEmail().GetValue()
		user.Email = &email
	}
	return user
}

func ToGetResponseFromService(user *model.User) *desc.GetResponse {
	role := desc.Role_USER
	if value, ok := desc.Role_value[user.Info.Role]; ok {
		role = desc.Role(value)
	}

	resp := &desc.GetResponse{
		Id:        user.ID,
		Name:      user.Info.Name,
		Email:     user.Info.Email,
		Role:      role,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
	if user.UpdatedAt.Valid {
		resp.UpdatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return resp
}
