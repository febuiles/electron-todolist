<script lang="ts">
  import { get } from 'svelte/store';
  import { onMount } from 'svelte';
  import Column from '../lib/Column.svelte';

  import { AppHost } from "../../src/config"
  import { todoStore } from '../stores/todostore';
  import { userStore } from '../stores/userstore';
  import type { ColumnType, Todo, Columnable, User } from '../lib/types';
  import { getTodolist, createTodolist } from '../lib/todolists';

  let user: User | null = null;

  onMount(async () => {
    try {
      user = await window.userAPI.getUser();
    } catch (err) {
      console.error('Failed to fetch user:', err);
    }
    if (!user) return;
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

  async function handleNewList() {
    if (!$userStore) return;

    console.log("Current", $userStore.lastUsedTodolistId);

    const newList = await createTodolist();

    userStore.update((user) => {
      if (!user) return null;
      return { ...user, lastUsedTodolistId: newList.id };
    });
    console.log($userStore.lastUsedTodolistId);

    getTodolist($userStore.lastUsedTodolistId);
  }

  async function handleShareList() {
    // TODO: Implement share functionality
    console.log('Share list clicked');
  }

  async function handleJoinList() {
    // TODO: Implement join functionality
    console.log('Join list clicked');
  }

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

      const response = await fetch(`${AppHost}/todos/`, {
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

    await fetch(`${AppHost}/todos/update`, {
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

<header class="app-header">
  <h1>pooltasks</h1>
  <div class="header-buttons">
    <button class="secondary-button" on:click={handleShareList}>
      <svg class="button-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="16" height="16">
        <path d="M16 5l-4-4-4 4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M12 21V7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M5 13v6a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2v-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      Share
    </button>
    <button class="secondary-button" on:click={handleJoinList}>
      <svg class="button-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="16" height="16"><path d="M12 12a5 5 0 1 0 0-10 5 5 0 0 0 0 10Zm0 2c-5 0-9 2.5-9 5v2h18v-2c0-2.5-4-5-9-5Z" fill="currentColor"/></svg>
      Join
    </button>
    <button class="primary-button" on:click={handleNewList}>
      <svg class="button-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="16" height="16"><path d="M12 5v14m7-7H5" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
      New List
    </button>
  </div>
</header>

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
