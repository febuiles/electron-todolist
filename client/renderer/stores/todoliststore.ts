import { writable } from 'svelte/store';

export const todolistStore = writable<Todolist | null>(null);
