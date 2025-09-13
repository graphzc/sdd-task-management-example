package auth

import (
	"context"

	"github.com/graphzc/sdd-task-management-example/internal/dto"
	"github.com/graphzc/sdd-task-management-example/internal/services/user"
)

type Handler interface {
	Register(ctx context.Context, req *dto.UserRegisterRequest) (*dto.MessageResponse, error)
	Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error)
}

type handler struct {
	userService user.Service
}

// @WireSet("Handler")
func New(userService user.Service) Handler {
	return &handler{
		userService: userService,
	}
}

func (h *handler) Register(ctx context.Context, req *dto.UserRegisterRequest) (*dto.MessageResponse, error) {
	serviceInput := user.UserRegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.userService.Register(ctx, &serviceInput); err != nil {
		return nil, err
	}

	return &dto.MessageResponse{
		Message: "User registered successfully",
	}, nil
}

func (h *handler) Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
	serviceInput := user.UserLoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	accessToken, err := h.userService.Login(ctx, &serviceInput)
	if err != nil {
		return nil, err
	}

	return &dto.UserLoginResponse{
		AccessToken: accessToken,
	}, nil
}
