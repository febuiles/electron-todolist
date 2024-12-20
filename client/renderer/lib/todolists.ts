import { get } from "svelte/store"

import { AppHost } from "../../src/config"

import { todolistStore } from "../stores/todoliststore"
import { userStore } from "../stores/userstore"
import { todoStore } from "../stores/todostore"

export async function createTodolist(): Promise<any> {
  const user = get(userStore)

  if (!user) {
    throw new Error('Failed to create todolist: Invalid user')
  }
  const todolist = await createRemote()
  user.lastUsedTodolistId = todolist.id
  userStore.set(user)
  window.electron.ipcRenderer.invoke('update-user', get(userStore))
  todoStore.set([])
  return todolist
}

async function createRemote(): Promise<any> {
  const user = get(userStore)

  if (!user || !user.id) {
    throw new Error('Failed to create todolist: Invalid user')
  }

  const response = await fetch(`${AppHost}/todolists/`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ user_id: user.id }),
  })

  return response.json()
}

export async function getTodolist(todolistId: string): Promise<any> {
  const user = get(userStore)

  if (!user || !user.id) {
    throw new Error('Failed to get todolist: Invalid user')
  }

  const response = await fetch(`${AppHost}/todolists/${todolistId}`, {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
  })

  if (!response.ok) {
    throw new Error(`Failed to get todolist: ${response.statusText}`)
  }

  return await response.json()
}
