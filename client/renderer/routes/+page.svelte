<script lang="ts">
  import { onMount } from 'svelte';
  import Column from '../lib/Column.svelte';
  import { writable } from 'svelte/store';

  let user;

  onMount(async () => {
    try {
      user = await window.userAPI.getUser();
    } catch (err) {
      console.error('Failed to fetch user:', err);
    }
  });

  // this state machines represents the possible transitions for a TODO
  const StateTransitions: Record<ColumnType, ColumnType[]> = {
    'todo': ['ongoing'],
    'ongoing': ['done', 'todo'],
    'done': ['ongoing']
  };

  const columns: Column[] = [
    { id: 'todo', title: 'TODO' },
    { id: 'ongoing', title: 'ONGOING' },
    { id: 'done', title: 'DONE' }
  ];

  const todoStore = writable([]);

  let draggedTodo: Todo | null = null;
  let draggedOverColumn: ColumnType | null = null;
  let newTodoTitle = ""

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

  async function fetchTodos() {
    const response = await fetch('http://localhost:8080/todos');
    const todos = await response.json();
    todoStore.set(todos);
  }


  async function addNewTodo(targetColumn: ColumnType, newTodoTitle: string): Promise<void> {
    if (newTodoTitle.trim()) {
      const newTodo = {
        title: newTodoTitle.trim(),
        user_id: user.id,
        column: targetColumn,
        lastUpdated: new Date().toLocaleString(),
      };

      const response = await fetch('http://localhost:8080/todos', {
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

  fetchTodos(); // Fetch the initial list of todos
</script>

<div class="board">
  {#each columns as column}
    <Column
      {column}
      todos={$todoStore}
      {draggedTodo}
      {draggedOverColumn}
      handleDragOver={handleDragOver}
      handleDragLeave={handleDragLeave}
      handleDrop={handleDrop}
      handleDragEnd={handleDragEnd}
      handleDragStart={handleDragStart}
      bind:newTodoTitle
      addNewTodo={addNewTodo}
    />
  {/each}
</div>
