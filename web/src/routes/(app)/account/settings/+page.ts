import type { PageLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageLoad = async ({ parent }) => {
	const { session } = await parent();

	// require user to be logged in
	if (!session) {
		throw redirect(303, '/login');
	}

	return {};
};
