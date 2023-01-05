import type { PageLoad } from './$types';
import { redirect } from '@sveltejs/kit';

const files = import.meta.glob('./*.md', {
	as: 'raw'
});

// Note: Redirect until docs are implemented
export const load: PageLoad = async (event) => {
	const { params } = event;
	const { file } = params;

	// dynamically load file
	const importFile = files[`./${file}.md`];

	// if file doesn't exist, redirect to /docs/welcome
	if (!importFile) throw redirect(303, '/docs/welcome');

	const raw = await importFile();

	return {
		raw
	};
};
