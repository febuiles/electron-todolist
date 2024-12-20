<script lang="ts">
  import { onMount } from 'svelte';
  import Column from '../lib/Column.svelte';
  import type { Columnable, User } from '../lib/types';
  import { createTodolist } from '../lib/todolists';
  import { getTodos } from '../lib/todos';
  import { userStore } from '../stores/userstore';
  import { todolistStore } from '../stores/todoliststore';

  onMount(async () => {
    try {
      let user = await window.userAPI.getUser();
      if (!user) return;

      userStore.set(user)
      getTodos(user.lastUsedTodolistId)
    } catch (err) {
      console.error('Failed to fetch user:', err);
    }
  });

  const columns: Columnable[] = [
    { id: 'todo', title: 'TODO' },
    { id: 'ongoing', title: 'ONGOING' },
    { id: 'done', title: 'DONE' }
  ];

  let newTodoTitle: string = ""
  let showJoinModal = false
  let joinSharedKey = ""
  let hostInput = ""
  let showShareModal = false
  const shareKey = "aws-msft-gcp"  // TODO fetch from the database

  async function handleNewList() {
    const newList = await createTodolist();
    todolistStore.set(newList)
    getTodos(newList.id);
  }

  async function handleShareList() {
    showShareModal = true;
  }

  async function handleJoinList() {
    showJoinModal = true;
  }

  function closeJoinModal() {
    showJoinModal = false;
    joinSharedKey = '';
    hostInput = '';
  }

  async function handleConnect() {
    // TODO: Implement connection logic
    console.log('Connecting with:', { joinSharedKey, hostInput });
    closeJoinModal();
  }

  function closeShareModal() {
    showShareModal = false;
  }

  async function copyToClipboard() {
    try {
      await navigator.clipboard.writeText(shareKey);
    } catch (err) {
      console.error('Failed to copy text:', err);
    }
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
      bind:newTodoTitle
    />
  {/each}
</div>

{#if showJoinModal}
  <div class="modal-overlay" on:click={closeJoinModal}>
    <div class="modal" on:click|stopPropagation>
      <h2>Join Tasklist</h2>
      <div class="modal-content">
        <div class="input-group">
          <label for="shared-key">Shared Key</label>
          <input
            id="shared-key"
            type="text"
            bind:value={joinSharedKey}
            placeholder="Enter shared key"
          />
        </div>
        <div class="input-group">
          <label for="host">Host (optional)</label>
          <input
            id="host"
            type="text"
            bind:value={hostInput}
            placeholder="Enter host address"
          />
        </div>
        <button class="primary-button big-button" on:click={handleConnect}>
          Connect
        </button>
      </div>
    </div>
  </div>
{/if}

{#if showShareModal}
  <div class="modal-overlay" on:click={closeShareModal}>
    <div class="modal" on:click|stopPropagation>
      <h2>Share Tasklist</h2>
      <div class="modal-content">
        <p class="share-instructions">Share this code with your friends to collaborate:</p>
        <div class="share-key-container">
          <div class="share-key">{shareKey}</div>
          <button class="copy-button" on:click={copyToClipboard} title="Copy to clipboard">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24">
              <path fill="none" d="M0 0h24v24H0z"/>
              <path d="M16 1H4c-1.1 0-2 .9-2 2v14h2V3h12V1zm3 4H8c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h11c1.1 0 2-.9 2-2V7c0-1.1-.9-2-2-2zm0 16H8V7h11v14z" fill="currentColor"/>
            </svg>
          </button>
        </div>
        <button class="primary-button big-button" on:click={closeShareModal}>
          Close
        </button>
      </div>
    </div>
  </div>
{/if}
