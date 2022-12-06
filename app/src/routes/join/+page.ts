import type { PageLoad } from './$types';
import { getSupabase } from '@supabase/auth-helpers-sveltekit';
import { redirect } from '@sveltejs/kit';

export const load: PageLoad = async (event) => {
	const { session } = await getSupabase(event);

	// require user to be logged in
	if (!session) {
		throw redirect(303, '/login');
	}

	// TODO: check if user is already off the waitlist or has a profile

	// do nothing
	return {};
};
