import type { PageLoad } from './$types';
import { getSupabase } from '@supabase/auth-helpers-sveltekit';
import { redirect } from '@sveltejs/kit';

//
export const load: PageLoad = async (event) => {
	const { session } = await getSupabase(event);

	// TODO: user is already logged in, send them elsewhere
	if (session) {
		throw redirect(303, '/');
	}

	// TODO: check if user is already off the waitlist or has a profile

	// do nothing
	return {};
};
