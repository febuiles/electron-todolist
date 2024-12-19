import { type Writable, writable } from 'svelte/store'
import type { Todo, User } from '../lib/types'

export const todoStore: Writable<Todo[]> = writable([])
