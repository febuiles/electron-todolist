import { userStore } from "../stores/userstore";
import { todoStore } from "../stores/todostore";
import { get } from "svelte/store";

export async function getTodolist(todolistId: number): Promise<void> {
  const response = await fetch(`http://localhost:8080/todolists/${todolistId}`);
  const todos = await response.json();
  todoStore.set(todos);
}

export async function createTodolist(): Promise<any> {
  const user = get(userStore);

  if (!user || !user.id) {
    throw new Error('Failed to create todolist: Invalid user');
  }

  const response = await fetch(`http://localhost:8080/todolists/`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ user_id: user.id }),
  });

  const todolist = await response.json();

  userStore.update((user) => {
    if (!user) return null;
    console.log({ ...user, lastUsedTodolistId: todolist.id });

    return { ...user, lastUsedTodolistId: todolist.id };
  });

  await window.electron.ipcRenderer.invoke('update-user', get(userStore));

  todoStore.set([]);

  return todolist;
}
