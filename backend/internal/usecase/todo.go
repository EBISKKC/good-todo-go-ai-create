package usecase

import (
	"context"
	"errors"
	"time"

	"good-todo-go/internal/domain/model"
	"good-todo-go/internal/domain/repository"
	"good-todo-go/internal/pkg"
	"good-todo-go/internal/usecase/input"
	"good-todo-go/internal/usecase/output"
)

var (
	ErrTodoNotFound     = errors.New("todo not found")
	ErrNotTodoOwner     = errors.New("not todo owner")
	ErrUnauthorized     = errors.New("unauthorized")
)

type ITodoInteractor interface {
	List(ctx context.Context, userID string) ([]*output.TodoOutput, error)
	ListPublic(ctx context.Context, tenantID string) ([]*output.TodoOutput, error)
	Create(ctx context.Context, userID, tenantID string, input *input.CreateTodoInput) (*output.TodoOutput, error)
	Update(ctx context.Context, userID string, input *input.UpdateTodoInput) (*output.TodoOutput, error)
	Delete(ctx context.Context, userID, todoID string) error
}

type TodoInteractor struct {
	todoRepo      repository.ITodoRepository
	uuidGenerator pkg.IUUIDGenerator
}

func NewTodoInteractor(todoRepo repository.ITodoRepository, uuidGenerator pkg.IUUIDGenerator) ITodoInteractor {
	return &TodoInteractor{
		todoRepo:      todoRepo,
		uuidGenerator: uuidGenerator,
	}
}

func (i *TodoInteractor) List(ctx context.Context, userID string) ([]*output.TodoOutput, error) {
	todos, err := i.todoRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*output.TodoOutput, len(todos))
	for idx, todo := range todos {
		result[idx] = toTodoOutput(todo)
	}
	return result, nil
}

func (i *TodoInteractor) ListPublic(ctx context.Context, tenantID string) ([]*output.TodoOutput, error) {
	todos, err := i.todoRepo.FindPublicByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	result := make([]*output.TodoOutput, len(todos))
	for idx, todo := range todos {
		result[idx] = toTodoOutput(todo)
	}
	return result, nil
}

func (i *TodoInteractor) Create(ctx context.Context, userID, tenantID string, inp *input.CreateTodoInput) (*output.TodoOutput, error) {
	todo := &model.Todo{
		ID:          i.uuidGenerator.Generate(),
		TenantID:    tenantID,
		UserID:      userID,
		Title:       inp.Title,
		Description: inp.Description,
		Completed:   false,
		IsPublic:    inp.IsPublic,
		DueDate:     inp.DueDate,
	}

	created, err := i.todoRepo.Create(ctx, todo)
	if err != nil {
		return nil, err
	}

	return toTodoOutput(created), nil
}

func (i *TodoInteractor) Update(ctx context.Context, userID string, inp *input.UpdateTodoInput) (*output.TodoOutput, error) {
	todo, err := i.todoRepo.FindByID(ctx, inp.ID)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, ErrTodoNotFound
	}

	if todo.UserID != userID {
		return nil, ErrNotTodoOwner
	}

	todo.Title = inp.Title
	todo.Description = inp.Description
	todo.Completed = inp.Completed
	todo.IsPublic = inp.IsPublic
	todo.DueDate = inp.DueDate

	if inp.Completed && todo.CompletedAt == nil {
		now := time.Now()
		todo.CompletedAt = &now
	} else if !inp.Completed {
		todo.CompletedAt = nil
	}

	updated, err := i.todoRepo.Update(ctx, todo)
	if err != nil {
		return nil, err
	}

	return toTodoOutput(updated), nil
}

func (i *TodoInteractor) Delete(ctx context.Context, userID, todoID string) error {
	todo, err := i.todoRepo.FindByID(ctx, todoID)
	if err != nil {
		return err
	}
	if todo == nil {
		return ErrTodoNotFound
	}

	if todo.UserID != userID {
		return ErrNotTodoOwner
	}

	return i.todoRepo.Delete(ctx, todoID)
}

func toTodoOutput(todo *model.Todo) *output.TodoOutput {
	return &output.TodoOutput{
		ID:          todo.ID,
		UserID:      todo.UserID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
		IsPublic:    todo.IsPublic,
		DueDate:     todo.DueDate,
		CompletedAt: todo.CompletedAt,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}
}
