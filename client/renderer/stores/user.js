import { writable } from 'svelte/store';

export const user = writable(null);

window.userAPI.getUser().then((u) => {
  user.set(u);
}).catch((err) => {
  console.error('Failed to fetch user:', err);
});
