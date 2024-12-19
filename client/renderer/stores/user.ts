import { type Writable, writable } from 'svelte/store';
import type { User } from '../lib/types';

export const userStore: Writable<User | null> = writable(null);

(window as any).userAPI.getUser().then((u: User) => {
  userStore.set(u);
}).catch((err: Error) => {
  console.error('Failed to fetch user:', err);
});
