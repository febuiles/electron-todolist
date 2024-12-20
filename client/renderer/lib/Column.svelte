<script lang="ts">
  import Card from './Card.svelte'
  import { todoStore } from '../stores/todostore.js'
  import { handleDragStart, handleDragOver, handleDragLeave, handleDrop, handleDragEnd } from '../lib/dragdrop'
  import { draggedTodo, draggedOverColumn } from '../stores/dragdropstore'
  import { createTodo } from './todos'

  export let column

  const StateTransitions: Record<string, string[]> = {
    'todo': ['ongoing'],
    'ongoing': ['done', 'todo'],
    'done': ['ongoing']
  };

  let newTodoTitle: string = ""

  function isValidTransition(fromState: string, toState: string) {
    return fromState === toState || StateTransitions[fromState]?.includes(toState)
  }
</script>

<div
  class="column"
  class:invalid-target={$draggedTodo && $draggedOverColumn === column.id && !isValidTransition($draggedTodo.column, column.id)}
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
        class="primary-button big-button"
        >
        + Add Item
      </button>
    </div>
  {/if}
</div>
