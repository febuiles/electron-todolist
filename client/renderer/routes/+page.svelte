<script lang="ts">
  import { get } from 'svelte/store';
  import { onMount } from 'svelte';
  import Column from '../lib/Column.svelte';
  import { todoStore } from '../stores/todostore.js';
  import { userStore } from '../stores/user.js';
  import type { ColumnType, Todo, Columnable, User } from '../lib/types';
  import { getTodolist, createTodolist } from '../lib/todolists';

  let user: User | null = null;

  onMount(async () => {
    try {
      user = await window.userAPI.getUser();
    } catch (err) {
      console.error('Failed to fetch user:', err);
    }
    userStore.set(user)
    getTodolist(user.lastUsedTodolistId)
  });

  // this state machines represents the possible transitions for a TODO
  const StateTransitions: Record<ColumnType, ColumnType[]> = {
    'todo': ['ongoing'],
    'ongoing': ['done', 'todo'],
    'done': ['ongoing']
  };

  const columns: Columnable[] = [
    { id: 'todo', title: 'TODO' },
    { id: 'ongoing', title: 'ONGOING' },
    { id: 'done', title: 'DONE' }
  ];

  let draggedTodo: Todo | null = null;
  let draggedOverColumn: ColumnType | null = null;
  let newTodoTitle: string = ""

  function handleDragStart(todo: Todo): void {
    draggedTodo = todo;
  }

  function handleDragOver(e: DragEvent, targetColumn: ColumnType): void {
    e.preventDefault();
    draggedOverColumn = targetColumn;

    // remove the mouse effect on invalid transitions
    if (draggedTodo && !isValidTransition(draggedTodo.column, targetColumn)) {
      e.dataTransfer!.dropEffect = 'none';
    }
  }

  function handleDragLeave(): void {
    draggedOverColumn = null;
  }

  function isValidTransition(fromState: ColumnType, toState: ColumnType): boolean {
    if (fromState === toState) {
      return true;
    }

    return StateTransitions[fromState]?.includes(toState);
  }

  function handleDrop(targetColumn: ColumnType): void {
    if (!draggedTodo) { return; }

    if (isValidTransition(draggedTodo.column, targetColumn)) {
      updateTodoColumn(draggedTodo.id, targetColumn);

      todoStore.update((todos) =>
        todos.map((todo) =>
          todo.id === draggedTodo!.id
            ? { ...todo, column: targetColumn, lastUpdated: new Date().toLocaleString() }
          : todo
        )
      );
    }

    draggedTodo = null;
    draggedOverColumn = null;
  }

  function handleDragEnd(): void {
    draggedTodo = null;
    draggedOverColumn = null;
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

      const response = await fetch('http://localhost:8080/todos/', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newTodo),
      });

      const createdTodo = await response.json();
      todoStore.update((todos) => [...todos, createdTodo]);
    }
  }


  async function updateTodoColumn(todoId: number, targetColumn: ColumnType): Promise<void> {
    const lastUpdated = new Date().toLocaleString();

    await fetch('http://localhost:8080/todos/update', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: todoId, column: targetColumn, lastUpdated }),
    });

    todoStore.update((todos) =>
      todos.map((todo) =>
        todo.id === todoId ? { ...todo, column: targetColumn, lastUpdated } : todo
      )
    );
  }
</script>

<div class="board">
  {#each columns as column}
    <Column
      {column}
      {draggedTodo}
      {draggedOverColumn}
      handleDragOver={handleDragOver}
      handleDragLeave={handleDragLeave}
      handleDrop={handleDrop}
      handleDragEnd={handleDragEnd}
      handleDragStart={handleDragStart}
      bind:newTodoTitle
      createTodo={createTodo}
    />
  {/each}
</div>
