<script lang="ts">
  import Card from './Card.svelte';
  import { todoStore } from '../stores/todostore.js';

  export let column;
  export let draggedTodo;
  export let draggedOverColumn;
  export let handleDragOver;
  export let handleDragLeave;
  export let handleDrop;
  export let handleDragEnd;
  export let handleDragStart;
  export let newTodoTitle: string;
  export let createTodo;

  function isValidTransition(fromState: string, toState: string) {
    const StateTransitions: Record<string, string[]> = {
      'todo': ['ongoing'],
      'ongoing': ['done', 'todo'],
      'done': ['ongoing']
    };
    return fromState === toState || StateTransitions[fromState]?.includes(toState);
  }
</script>

<div
  class="column"
  class:invalid-target={draggedTodo && draggedOverColumn === column.id && !isValidTransition(draggedTodo.column, column.id)}
  on:dragover={(e) => handleDragOver(e, column.id)}
  on:dragleave={handleDragLeave}
  on:drop={() => handleDrop(column.id)}
  on:dragend={handleDragEnd}
>
  <div class="column-header">
    <h2>{column.title}</h2>
    <span class="task-count">
      {$todoStore.filter(todo => todo.column === column.id)?.length || 0}
    </span>
  </div>

  <div class="cards">
    {#if $todoStore.filter(todo => todo.column === column.id).length > 0}
      {#each $todoStore.filter(todo => todo.column === column.id) as todo}
        <Card {todo} {handleDragStart}/>
      {/each}
    {/if}
  </div>

  {#if column.id === 'todo'}
    <div class="add-card-section">
      <input
        type="text"
        placeholder="Enter new TODO..."
        bind:value={newTodoTitle}
        class="add-card-input"
        on:keydown={(e) => {
          if (e.key === 'Enter') {
            createTodo(column.id, newTodoTitle);
            newTodoTitle = "";
          }
        }}
      />
      <button
        on:click={() => {
          createTodo(column.id, newTodoTitle);
          newTodoTitle = "";
        }}
        class="primary-button add-card-button"
        >
        + Add Item
      </button>
    </div>
  {/if}
</div>
