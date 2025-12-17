package router

import (
	"good-todo-go/internal/presentation/public/api"

	"github.com/labstack/echo/v4"
)

func (s *Server) ListTodos(ctx echo.Context) error {
	return s.todoController.ListTodos(ctx)
}

func (s *Server) ListPublicTodos(ctx echo.Context) error {
	return s.todoController.ListPublicTodos(ctx)
}

func (s *Server) CreateTodo(ctx echo.Context) error {
	var req api.CreateTodoRequest
	if err := ctx.Bind(&req); err != nil {
		return err
	}
	return s.todoController.CreateTodo(ctx, req)
}

func (s *Server) UpdateTodo(ctx echo.Context, id string) error {
	var req api.UpdateTodoRequest
	if err := ctx.Bind(&req); err != nil {
		return err
	}
	return s.todoController.UpdateTodo(ctx, id, req)
}

func (s *Server) DeleteTodo(ctx echo.Context, id string) error {
	return s.todoController.DeleteTodo(ctx, id)
}
