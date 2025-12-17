package repository

import (
	"context"
	"fmt"

	"good-todo-go/internal/domain/model"
	"good-todo-go/internal/domain/repository"
	"good-todo-go/internal/ent/generated"
	"good-todo-go/internal/ent/generated/todo"
)

type TodoRepository struct {
	client *generated.Client
}

func NewTodoRepository(client *generated.Client) repository.ITodoRepository {
	return &TodoRepository{client: client}
}

func (r *TodoRepository) Create(ctx context.Context, t *model.Todo) (*model.Todo, error) {
	builder := r.client.Todo.Create().
		SetID(t.ID).
		SetTenantID(t.TenantID).
		SetUserID(t.UserID).
		SetTitle(t.Title).
		SetDescription(t.Description).
		SetCompleted(t.Completed).
		SetIsPublic(t.IsPublic)

	if t.DueDate != nil {
		builder.SetDueDate(*t.DueDate)
	}
	if t.CompletedAt != nil {
		builder.SetCompletedAt(*t.CompletedAt)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}
	return toModelTodo(created), nil
}

func (r *TodoRepository) FindByID(ctx context.Context, id string) (*model.Todo, error) {
	t, err := r.client.Todo.Get(ctx, id)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find todo by id: %w", err)
	}
	return toModelTodo(t), nil
}

func (r *TodoRepository) FindByUserID(ctx context.Context, userID string) ([]*model.Todo, error) {
	todos, err := r.client.Todo.Query().
		Where(todo.UserIDEQ(userID)).
		Order(generated.Desc(todo.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find todos by user id: %w", err)
	}

	result := make([]*model.Todo, len(todos))
	for i, t := range todos {
		result[i] = toModelTodo(t)
	}
	return result, nil
}

func (r *TodoRepository) FindPublicByTenantID(ctx context.Context, tenantID string) ([]*model.Todo, error) {
	todos, err := r.client.Todo.Query().
		Where(
			todo.TenantIDEQ(tenantID),
			todo.IsPublicEQ(true),
		).
		Order(generated.Desc(todo.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find public todos by tenant id: %w", err)
	}

	result := make([]*model.Todo, len(todos))
	for i, t := range todos {
		result[i] = toModelTodo(t)
	}
	return result, nil
}

func (r *TodoRepository) Update(ctx context.Context, t *model.Todo) (*model.Todo, error) {
	builder := r.client.Todo.UpdateOneID(t.ID).
		SetTitle(t.Title).
		SetDescription(t.Description).
		SetCompleted(t.Completed).
		SetIsPublic(t.IsPublic)

	if t.DueDate != nil {
		builder.SetDueDate(*t.DueDate)
	} else {
		builder.ClearDueDate()
	}
	if t.CompletedAt != nil {
		builder.SetCompletedAt(*t.CompletedAt)
	} else {
		builder.ClearCompletedAt()
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}
	return toModelTodo(updated), nil
}

func (r *TodoRepository) Delete(ctx context.Context, id string) error {
	err := r.client.Todo.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}
	return nil
}

func toModelTodo(t *generated.Todo) *model.Todo {
	return &model.Todo{
		ID:          t.ID,
		TenantID:    t.TenantID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		IsPublic:    t.IsPublic,
		DueDate:     t.DueDate,
		CompletedAt: t.CompletedAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
