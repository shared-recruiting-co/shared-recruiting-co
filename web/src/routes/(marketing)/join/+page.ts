import type { PageLoad } from './$types';
import { redirect } from '@sveltejs/kit';

type Data = {
	success: boolean;
};

export const load: PageLoad<Data> = async ({ parent }) => {
	const { session, supabase } = await parent();
	// require user to be logged in
	if (!session) {
		throw redirect(303, '/login');
	}

	// if user is already on the waitlist,
	const { data: waitlist } = await supabase.from('waitlist').select('*').maybeSingle();
	if (waitlist) {
		// if they can create an account, send them to the profile page
		if (waitlist.can_create_account) {
			throw redirect(303, '/account/setup');
		}
		// if they cannot, show them the success state
		return {
			success: true
		};
	}

	// show form
	return {
		success: false
	};
};
