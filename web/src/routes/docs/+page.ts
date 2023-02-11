import type { PageLoad } from './$types';
import { redirect } from '@sveltejs/kit';

// Redirect to welcome
export const load: PageLoad = async (event) => {
	throw redirect(303, '/docs/welcome');
};
