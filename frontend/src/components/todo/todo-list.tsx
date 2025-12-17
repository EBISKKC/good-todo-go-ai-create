'use client';

import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import axiosInstance from '@/api/axios-instance';
import { TodoItem } from './todo-item';
import { CreateTodoDialog } from './create-todo-dialog';
import { EditTodoDialog } from './edit-todo-dialog';
import { Button } from '@/components/ui/button';
import { Plus } from 'lucide-react';
import { TodoResponse, TodoListResponse } from '@/api/public/model';
import { toast } from 'sonner';

export function TodoList() {
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [editingTodo, setEditingTodo] = useState<TodoResponse | null>(null);
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery<TodoListResponse>({
    queryKey: ['todos'],
    queryFn: async () => {
      const response = await axiosInstance.get('/todos');
      return response.data;
    },
  });

  const updateMutation = useMutation({
    mutationFn: async ({
      id,
      data,
    }: {
      id: string;
      data: { title: string; description?: string; completed: boolean; is_public?: boolean; due_date?: string | null };
    }) => {
      const response = await axiosInstance.put(`/todos/${id}`, data);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['todos'] });
      toast.success('Todo updated');
    },
    onError: () => {
      toast.error('Failed to update todo');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await axiosInstance.delete(`/todos/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['todos'] });
      toast.success('Todo deleted');
    },
    onError: () => {
      toast.error('Failed to delete todo');
    },
  });

  const handleToggle = (id: string, completed: boolean) => {
    const todo = data?.todos.find((t) => t.id === id);
    if (todo) {
      updateMutation.mutate({
        id,
        data: {
          title: todo.title,
          description: todo.description,
          completed,
          is_public: todo.is_public,
          due_date: todo.due_date,
        },
      });
    }
  };

  const handleDelete = (id: string) => {
    if (confirm('Are you sure you want to delete this todo?')) {
      deleteMutation.mutate(id);
    }
  };

  if (isLoading) {
    return <div className="text-center py-8">Loading...</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold">My Todos</h2>
        <Button onClick={() => setIsCreateOpen(true)}>
          <Plus className="h-4 w-4 mr-2" />
          Add Todo
        </Button>
      </div>

      {data?.todos.length === 0 ? (
        <div className="text-center py-8 text-gray-500">
          No todos yet. Create your first todo!
        </div>
      ) : (
        <div className="space-y-2">
          {data?.todos.map((todo) => (
            <TodoItem
              key={todo.id}
              todo={todo}
              onToggle={handleToggle}
              onEdit={setEditingTodo}
              onDelete={handleDelete}
            />
          ))}
        </div>
      )}

      <CreateTodoDialog open={isCreateOpen} onOpenChange={setIsCreateOpen} />
      {editingTodo && (
        <EditTodoDialog
          todo={editingTodo}
          open={!!editingTodo}
          onOpenChange={(open) => !open && setEditingTodo(null)}
        />
      )}
    </div>
  );
}
