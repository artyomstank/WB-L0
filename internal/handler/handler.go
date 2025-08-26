package handler

import "L0-wb/internal/service"

type UserHandler struct {
	UserSevice service.UserService
}

func NewHandler(us service.UserService) UserHandler {
	return UserHandler{
		UserSevice: us,
	}
}
