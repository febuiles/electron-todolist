import { get } from 'svelte/store';
import { todoStore } from '../stores/todostore';
import { userStore } from '../stores/userstore';
import type { ColumnType } from './types';
import { AppHost } from '../../src/config';

export async function getTodos(todolistId: number): Promise<void> {
  const response = await fetch(`${AppHost}/todolists/${todolistId}`);
  const todos = await response.json();
  todoStore.set(todos);
}

export async function createTodo(targetColumn: ColumnType, newTodoTitle: string): Promise<void> {
  const user = get(userStore);

  if (!user || !user.id) {
    throw new Error('Failed to add new todo: Invalid user')
  }

  if (newTodoTitle.trim()) {
    const newTodo = {
      title: newTodoTitle.trim(),
      user_id: user.id,
      column: targetColumn,
      lastUpdated: new Date().toLocaleString(),
      todolist_id: user.lastUsedTodolistId
    };

    const response = await fetch(`${AppHost}/todos/`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newTodo),
    });

    const createdTodo = await response.json();
    todoStore.update((todos) => [...todos, createdTodo]);
  }
}

export async function updateTodo(todoId: number, targetColumn: ColumnType): Promise<void> {
  const lastUpdated = new Date().toLocaleString();

  await fetch(`${AppHost}/todos/${todoId}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ id: todoId, column: targetColumn, lastUpdated }),
  });

  todoStore.update((todos) =>
    todos.map((todo) =>
      todo.id === todoId ? { ...todo, column: targetColumn, lastUpdated } : todo
             )
                  );
}
