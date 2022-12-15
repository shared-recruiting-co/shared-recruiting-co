import type { PageLoad } from './$types';
import { getSupabase } from '@supabase/auth-helpers-sveltekit';
import { redirect } from '@sveltejs/kit';

export const load: PageLoad = async (event) => {
	const { session, supabaseClient } = await getSupabase(event);

	if (!session) return {};

	const [{ data: profile }, { data: waitlist }] = await Promise.all([
		supabaseClient.from('user_profile').select('*').maybeSingle(),
		supabaseClient.from('waitlist').select('*').maybeSingle()
	]);
	if (profile) {
		throw redirect(303, '/account/profile');
	}

	// if user is already on the waitlist,
	if (waitlist && waitlist.can_create_account) {
		// if they can create an account, send them to create one
		throw redirect(303, '/account/profile/create');
	}

	// send to waitlist page
	throw redirect(303, '/join');
};
