import type { PageLoad } from './$types';
// import { getSupabase } from '@supabase/auth-helpers-sveltekit';
// import { redirect } from '@sveltejs/kit';

export const load: PageLoad = async (event) => {
	// const { session } = await getSupabase(event);

	// TODO: Alternatively, I can show a login button
//
	// // require user to be logged in
	// if (!session) {
		// throw redirect(303, '/');
	// }

	// TODO: check if user is already off the waitlist

	// do nothing
	return {};
};
