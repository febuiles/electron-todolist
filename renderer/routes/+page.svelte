<script lang="ts">
  import { writable } from 'svelte/store';

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

  const todoStore: Writable<Todo[]> = writable([
    {
      id: 1,
      title: 'Implement authentication',
      creator: 'John Doe',
      column: 'todo'
    },
    {
      id: 2,
      title: 'Create database schema',
      creator: 'Jane Smith',
      column: 'todo'
    }
  ]);

  let draggedCard: Todo | null = null;
  let draggedOverColumn: ColumnType | null = null;
  let nextId = 3;
  let newCardTitle = ""

  function handleDragStart(card: Todo): void {
    draggedCard = card;
  }

  function handleDragOver(e: DragEvent, targetColumn: ColumnType): void {
    e.preventDefault();
    draggedOverColumn = targetColumn;

    // remove the mouse effect on invalid transitions
    if (draggedCard && !isValidTransition(draggedCard.column, targetColumn)) {
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
    if (!draggedCard) {
      return;
    }

    if (isValidTransition(draggedCard.column, targetColumn)) {
      todoStore.update(todos =>
        todos.map(todo =>
          todo.id === draggedCard!.id
            ? { ...todo, column: targetColumn }
          : todo
        )
      );
    }
    draggedCard = null;
    draggedOverColumn = null;
  }

  function handleDragEnd(): void {
    draggedCard = null;
    draggedOverColumn = null;
  }

  function addNewCard(targetColumn: ColumnType): void {
    if (newCardTitle.trim()) {
      todoStore.update(todos => [
        ...todos,
        {
          id: nextId++,
          title: newCardTitle.trim(),
          creator: "Anonymous",
          column: targetColumn
        }
      ]);
      newCardTitle = "";
    }
  }
</script>

<div class="board">
  {#each columns as column}
    <div
      class="column"
      class:invalid-target={draggedCard && draggedOverColumn === column.id && !isValidTransition(draggedCard.column, column.id)}
      on:dragover={(e) => handleDragOver(e, column.id)}
      on:dragleave={handleDragLeave}
      on:drop={() => handleDrop(column.id)}
      on:dragend={handleDragEnd}
      >
      <div class="column-header">
        <h2>{column.title}</h2>
        <span class="task-count">
          {$todoStore.filter(card => card.column === column.id).length}
        </span>
      </div>

      <div class="cards">
        {#if $todoStore.filter(card => card.column === column.id).length > 0}
          {#each $todoStore.filter(card => card.column === column.id) as card}
            <div
              class="card"
              draggable="true"
              on:dragstart={() => handleDragStart(card)}
              >
              <h3>{card.title}</h3>
              <div class="card-meta">
                <span>Created by: {card.creator}</span>
              </div>
            </div>
          {/each}
        {/if}
      </div>

      {#if column.id === 'todo'}
        <div class="add-card-section">
          <input
            type="text"
            placeholder="Enter new TODO..."
            bind:value={newCardTitle}
            class="add-card-input"
            on:keydown={(e) => { if (e.key === 'Enter') addNewCard(column.id); }}
          />
          <button on:click={() => addNewCard(column.id)} class="add-card-button">
            + Add Item
          </button>
        </div>
      {/if}
    </div>
  {/each}
</div>
