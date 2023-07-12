/** @type {import('tailwindcss').Config} */
module.exports = {
	content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
    darkMode: 'class',
	theme: {
		extend: {
            fontFamily: {
                titulo: ['"Secular One"', '"sans-serif"'],
                texto: ['"Oxygen"', '"sans-serif"'],
            }
        },
	},
	plugins: [],
}
