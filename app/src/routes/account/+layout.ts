import type { PageLoad } from './$types';
import { getSupabase } from '@supabase/auth-helpers-sveltekit';
import { redirect } from '@sveltejs/kit';


export const load: PageLoad = async (event) => {
	const { session, supabaseClient } = await getSupabase(event);

	// require user to be logged in
	if (!session) {
		throw redirect(303, '/login');
	}

	// if user has a profile, we are good
	const { data: profile } = await supabaseClient.from('user_profile').select('*').maybeSingle();
	if (profile) {
		return {}
	}

	// if the are aren't on the waitlist or can't create an account, redirect to the waitlist page
	const { data: waitlist } = await supabaseClient.from('waitlist').select('*').maybeSingle();
	if (!waitlist || !waitlist.can_create_account) {
		throw redirect(303, '/join');
	}

	// send to profile creation page
	throw redirect(303, '/account/profile');
};
