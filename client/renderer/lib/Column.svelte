<script lang="ts">
  import Card from './Card.svelte';

  export let column;
  export let todos;
  export let todoStore;
  export let draggedTodo;
  export let draggedOverColumn;
  export let handleDragOver;
  export let handleDragLeave;
  export let handleDrop;
  export let handleDragEnd;
  export let handleDragStart;
  export let newTodoTitle;
  export let addNewTodo;

  function isValidTransition(fromState, toState) {
    const StateTransitions = {
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
      {todos.filter(todo => todo.column === column.id).length}
    </span>
  </div>

  <div class="cards">
    {#if todos.filter(todo => todo.column === column.id).length > 0}
      {#each todos.filter(todo => todo.column === column.id) as todo}
        <Card {todo} {todoStore} {handleDragStart}/>
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
        on:keydown={(e) => { if (e.key === 'Enter') addNewTodo(column.id, newTodoTitle); }}
      />
      <button on:click={() => addNewTodo(column.id, newTodoTitle)} class="add-card-button">
        + Add Item
      </button>
    </div>
  {/if}
</div>
