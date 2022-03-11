const config = {
	content: ['./src/**/*.{html,js,svelte,ts}'],

	darkMode: 'media',

	theme: {
		extend: {}
	},

	plugins: [require('@tailwindcss/forms')]
};

module.exports = config;
