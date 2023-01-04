import type { PageLoad } from './$types';
import { redirect } from '@sveltejs/kit';

// Note: Redirect until docs are implemented
export const load: PageLoad = async (event) => {
	throw redirect(303, '/');
};
