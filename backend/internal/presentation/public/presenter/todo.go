package presenter

import (
	"net/http"

	"good-todo-go/internal/presentation/public/api"
	"good-todo-go/internal/usecase/output"

	"github.com/labstack/echo/v4"
)

type ITodoPresenter interface {
	List(c echo.Context, todos []*output.TodoOutput) error
	Create(c echo.Context, todo *output.TodoOutput) error
	Update(c echo.Context, todo *output.TodoOutput) error
	Delete(c echo.Context) error
}

type TodoPresenter struct{}

func NewTodoPresenter() ITodoPresenter {
	return &TodoPresenter{}
}

func (p *TodoPresenter) List(c echo.Context, todos []*output.TodoOutput) error {
	result := make([]api.TodoResponse, len(todos))
	for i, todo := range todos {
		result[i] = toTodoResponse(todo)
	}
	return c.JSON(http.StatusOK, api.TodoListResponse{Todos: result})
}

func (p *TodoPresenter) Create(c echo.Context, todo *output.TodoOutput) error {
	return c.JSON(http.StatusCreated, toTodoResponse(todo))
}

func (p *TodoPresenter) Update(c echo.Context, todo *output.TodoOutput) error {
	return c.JSON(http.StatusOK, toTodoResponse(todo))
}

func (p *TodoPresenter) Delete(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func toTodoResponse(out *output.TodoOutput) api.TodoResponse {
	resp := api.TodoResponse{
		Id:          out.ID,
		UserId:      out.UserID,
		Title:       out.Title,
		Description: out.Description,
		Completed:   out.Completed,
		IsPublic:    out.IsPublic,
		CreatedAt:   out.CreatedAt,
		UpdatedAt:   out.UpdatedAt,
	}
	if out.DueDate != nil {
		resp.DueDate = out.DueDate
	}
	if out.CompletedAt != nil {
		resp.CompletedAt = out.CompletedAt
	}
	return resp
}
