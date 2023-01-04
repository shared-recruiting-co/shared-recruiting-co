import type { PageLoad } from './$types';
import { getSupabase } from '@supabase/auth-helpers-sveltekit';
import { redirect } from '@sveltejs/kit';

export const load: PageLoad = async (event) => {
	const { session, supabaseClient } = await getSupabase(event);
	const { route } = event;

	// require user to be logged in
	if (!session) {
		throw redirect(303, '/login');
	}

	const isAccountCreationPage =
		route.id === '/account/profile/create' || route.id === '/account/profile/connect';

	// if user has a profile, we are good
	const { data: profile } = await supabaseClient.from('user_profile').select('*').maybeSingle();
	if (profile) {
		if (isAccountCreationPage) {
			throw redirect(303, '/account/profile');
		}

		return {
			profile: {
				firstName: profile.first_name,
				lastName: profile.last_name,
				email: session.user.email,
				createdAt: profile.created_at,
				isActive: profile.is_active,
				autoArchive: profile.auto_archive,
				autoContribute: profile.auto_contribute
			}
		};
	}

	// if the are aren't on the waitlist or can't create an account, redirect to the waitlist page
	const { data: waitlist } = await supabaseClient.from('waitlist').select('*').maybeSingle();
	if (!waitlist || !waitlist.can_create_account) {
		throw redirect(303, '/join');
	}

	// send to profile creation page
	if (!isAccountCreationPage) {
		throw redirect(303, '/account/profile/create');
	}

	if (waitlist.can_create_account) {
		return {
			profile: {
				firstName: waitlist.first_name,
				lastName: waitlist.last_name,
				email: session.user.email
			}
		};
	}

	return {};
};
