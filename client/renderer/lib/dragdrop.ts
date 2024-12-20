import { get } from 'svelte/store';
import type { ColumnType, Todo } from '../lib/types';
import { todoStore } from '../stores/todostore';
import { draggedTodo, draggedOverColumn } from '../stores/dragdropstore';
import { updateTodo } from './todos';

// represents the possible transitions for a TODO
const StateTransitions: Record<ColumnType, ColumnType[]> = {
  'todo': ['ongoing'],
  'ongoing': ['done', 'todo'],
  'done': ['ongoing']
}

export function handleDragStart(todo: Todo): void {
  draggedTodo.set(todo)
}

export function handleDragOver(e: DragEvent, targetColumn: ColumnType): void {
  e.preventDefault()
  draggedOverColumn.set(targetColumn);

  // remove the mouse effect on invalid transitions
  if (get(draggedTodo) && !isValidTransition(get(draggedTodo).column, targetColumn)) {
    e.dataTransfer!.dropEffect = 'none';
  }
}

export function handleDragLeave(): void {
  draggedOverColumn.set(null);
}

export function handleDrop(targetColumn: ColumnType): void {
  if (!draggedTodo) { return }
  let item = get(draggedTodo)

  if (isValidTransition(item.column, targetColumn)) {
    updateTodo(item.id, targetColumn)

    todoStore.update((todos) =>
      todos.map((todo) =>
        todo.id === draggedTodo!.id
        ? { ...todo, column: targetColumn, lastUpdated: new Date().toLocaleString() }
        : todo
      )
    )
  }

  draggedTodo.set(null)
  draggedOverColumn.set(null)
}

export function handleDragEnd(): void {
  draggedTodo.set(null)
  draggedOverColumn.set(null)
}

export function isValidTransition(fromState: ColumnType, toState: ColumnType): boolean {
  if (fromState === toState) {
    return true
  }

  return StateTransitions[fromState]?.includes(toState);
}
