package controller

import (
	"net/http"

	"good-todo-go/internal/presentation/public/api"
	"good-todo-go/internal/presentation/public/presenter"
	"good-todo-go/internal/presentation/public/router/context_keys"
	"good-todo-go/internal/usecase"
	"good-todo-go/internal/usecase/input"

	"github.com/labstack/echo/v4"
)

type TodoController struct {
	todoUsecase   usecase.ITodoInteractor
	todoPresenter presenter.ITodoPresenter
}

func NewTodoController(
	todoUsecase usecase.ITodoInteractor,
	todoPresenter presenter.ITodoPresenter,
) *TodoController {
	return &TodoController{
		todoUsecase:   todoUsecase,
		todoPresenter: todoPresenter,
	}
}

func (ctrl *TodoController) ListTodos(c echo.Context) error {
	userID, ok := c.Get(context_keys.UserIDContextKey).(string)
	if !ok || userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	todos, err := ctrl.todoUsecase.List(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctrl.todoPresenter.List(c, todos)
}

func (ctrl *TodoController) ListPublicTodos(c echo.Context) error {
	tenantID, ok := c.Get(context_keys.TenantIDContextKey).(string)
	if !ok || tenantID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	todos, err := ctrl.todoUsecase.ListPublic(c.Request().Context(), tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctrl.todoPresenter.List(c, todos)
}

func (ctrl *TodoController) CreateTodo(c echo.Context, req api.CreateTodoRequest) error {
	userID, ok := c.Get(context_keys.UserIDContextKey).(string)
	if !ok || userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	tenantID, ok := c.Get(context_keys.TenantIDContextKey).(string)
	if !ok || tenantID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	description := ""
	if req.Description != nil {
		description = *req.Description
	}
	isPublic := false
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	todo, err := ctrl.todoUsecase.Create(c.Request().Context(), userID, tenantID, &input.CreateTodoInput{
		Title:       req.Title,
		Description: description,
		IsPublic:    isPublic,
		DueDate:     req.DueDate,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctrl.todoPresenter.Create(c, todo)
}

func (ctrl *TodoController) UpdateTodo(c echo.Context, id string, req api.UpdateTodoRequest) error {
	userID, ok := c.Get(context_keys.UserIDContextKey).(string)
	if !ok || userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	description := ""
	if req.Description != nil {
		description = *req.Description
	}
	isPublic := false
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	todo, err := ctrl.todoUsecase.Update(c.Request().Context(), userID, &input.UpdateTodoInput{
		ID:          id,
		Title:       req.Title,
		Description: description,
		Completed:   req.Completed,
		IsPublic:    isPublic,
		DueDate:     req.DueDate,
	})
	if err != nil {
		if err == usecase.ErrTodoNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "todo not found")
		}
		if err == usecase.ErrNotTodoOwner {
			return echo.NewHTTPError(http.StatusForbidden, "not authorized")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctrl.todoPresenter.Update(c, todo)
}

func (ctrl *TodoController) DeleteTodo(c echo.Context, id string) error {
	userID, ok := c.Get(context_keys.UserIDContextKey).(string)
	if !ok || userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	err := ctrl.todoUsecase.Delete(c.Request().Context(), userID, id)
	if err != nil {
		if err == usecase.ErrTodoNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "todo not found")
		}
		if err == usecase.ErrNotTodoOwner {
			return echo.NewHTTPError(http.StatusForbidden, "not authorized")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctrl.todoPresenter.Delete(c)
}
