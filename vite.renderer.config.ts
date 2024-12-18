import { defineConfig } from 'vite';
import { sveltekit } from '@sveltejs/kit/vite';

//https://vitejs.dev/config
export default defineConfig({
    server: {
	fs: {
	    allow: [ './renderer' ],
	}
    },
    plugins: [sveltekit()],
    clearScreen: false,
});
