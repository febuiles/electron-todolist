import { writable } from 'svelte/store';

export const draggedTodo = writable<Todo | null>(null);
export const draggedOverColumn = writable<ColumnType | null>(null);
