/** @type {import('tailwindcss').Config} */
module.exports = {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			fontFamily: {
				sans: ['Miriam Libre', 'sans-serif']
			}
		}
	},
	plugins: [require("@tailwindcss/forms"), require('@tailwindcss/typography')]
};
