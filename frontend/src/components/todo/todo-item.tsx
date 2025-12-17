'use client';

import { Checkbox } from '@/components/ui/checkbox';
import { Button } from '@/components/ui/button';
import { Pencil, Trash2 } from 'lucide-react';
import { TodoResponse } from '@/api/public/model';

interface TodoItemProps {
  todo: TodoResponse;
  onToggle: (id: string, completed: boolean) => void;
  onEdit: (todo: TodoResponse) => void;
  onDelete: (id: string) => void;
}

export function TodoItem({ todo, onToggle, onEdit, onDelete }: TodoItemProps) {
  return (
    <div className="flex items-center gap-4 p-4 border rounded-lg hover:bg-gray-50">
      <Checkbox
        checked={todo.completed}
        onCheckedChange={(checked) => onToggle(todo.id, checked as boolean)}
      />
      <div className="flex-1 min-w-0">
        <h3
          className={`font-medium truncate ${
            todo.completed ? 'line-through text-gray-400' : ''
          }`}
        >
          {todo.title}
        </h3>
        {todo.description && (
          <p className="text-sm text-gray-500 truncate">{todo.description}</p>
        )}
        <div className="flex gap-2 mt-1 text-xs text-gray-400">
          {todo.is_public && <span className="px-2 py-0.5 bg-blue-100 text-blue-700 rounded">Public</span>}
          {todo.due_date && (
            <span>Due: {new Date(todo.due_date).toLocaleDateString()}</span>
          )}
        </div>
      </div>
      <div className="flex gap-2">
        <Button variant="ghost" size="icon" onClick={() => onEdit(todo)}>
          <Pencil className="h-4 w-4" />
        </Button>
        <Button variant="ghost" size="icon" onClick={() => onDelete(todo.id)}>
          <Trash2 className="h-4 w-4 text-red-500" />
        </Button>
      </div>
    </div>
  );
}
