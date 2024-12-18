import type { Config } from 'tailwindcss';

const config: Config = {
    content: ['./renderer/**/*.{html,js,svelte,ts}'],
    theme: {
	extend: {},
    },
    plugins: [],
}

export default config;
