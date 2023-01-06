import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

const files = import.meta.glob('./*.md', {
	as: 'raw'
});

type Data = {
	markdown: {
		source: string;
	};
};

// Note: Redirect until docs are implemented
export const load: PageServerLoad<Data> = async ({ params }) => {
	const { file } = params;

	// dynamically load file
	const importFile = files[`./${file}.md`];

	// if file doesn't exist, redirect to /docs/welcome
	if (!importFile) throw redirect(303, '/docs/welcome');

	const source = await importFile();

	return {
		markdown: {
			source
		}
	};
};

export const prerender = true;
