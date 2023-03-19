import type { PageLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageLoad = async ({ parent }) => {
	const { session, supabase } = await parent();
	if (!session) return {};

	const [{ data: profile }, { data: waitlist }] = await Promise.all([
		supabase.from('user_profile').select('*').maybeSingle(),
		supabase.from('waitlist').select('*').maybeSingle()
	]);

	if (profile) {
		throw redirect(303, '/account/profile');
	}

	// if user is already on the waitlist,
	if (waitlist && waitlist.can_create_account) {
		// if they can create an account, send them to create one
		throw redirect(303, '/account/setup');
	}

	// send to waitlist page
	throw redirect(303, '/join');
};
