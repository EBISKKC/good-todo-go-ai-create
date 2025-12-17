package router

import (
	"good-todo-go/internal/presentation/public/controller"
)

type Server struct {
	authController *controller.AuthController
	userController *controller.UserController
	todoController *controller.TodoController
}

func NewServer(
	authController *controller.AuthController,
	userController *controller.UserController,
	todoController *controller.TodoController,
) *Server {
	return &Server{
		authController: authController,
		userController: userController,
		todoController: todoController,
	}
}
