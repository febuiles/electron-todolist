<script lang="ts">
  import { AppHost } from "../../src/config"
  import type { Todo } from '../lib/types';
  import { todoStore } from '../stores/todostore.js';

  export let todo: Todo;
  export let handleDragStart;

  async function deleteTodo() {
    if (confirm(`Are you sure you want to delete "${todo.title}"?`)) {
      await handleDelete(todo.id);
    }
  }

  async function handleDelete(todoID: number) {
    try {
      const response = await fetch(`${AppHost}/todos/${todoID}`, {
        method: 'DELETE'
      });

      if (!response.ok) {
        throw new Error('Failed to delete the TODO');
      }
      todoStore.update(todos => todos.filter(todo => todo.id !== todoID));
    } catch (error: any) {
      alert(`Error: ${error.message}`);
    }
  }

</script>

<div
  class="card"
  draggable="true"
  on:dragstart={() => handleDragStart(todo)}
>
  <div class="card-header">
    <h3>{todo.title}</h3>
    <button class="delete-icon" on:click={deleteTodo} title="Delete TODO">
      &times;
    </button>
  </div>
  <div class="card-meta">
    <span>Created by: {todo.creator}</span>
    <span>Last updated: {todo.lastUpdated}</span>
  </div>
</div>
