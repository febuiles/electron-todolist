<script lang="ts">
  export let todo;
  export let handleDragStart;
</script>

<div
  class="card"
  draggable="true"
  on:dragstart={() => handleDragStart(todo)}
>
  <h3>{todo.title}</h3>
  <div class="card-meta">
    <span>Created by: {todo.creator}</span>
    <span>Last updated: {todo.lastUpdated}</span>
  </div>
</div>
