import type { PageLoad } from './$types';
import { getSupabase } from '@supabase/auth-helpers-sveltekit';
import { redirect } from '@sveltejs/kit';

export const load: PageLoad = async (event) => {
	const { session, supabaseClient } = await getSupabase(event);
	const path = event.route.id;

	// redirect to login
	if (!session) {
		if (path !== '/(app)/recruiter/login') {
			throw redirect(303, '/recruiter/login');
		}
		return {};
	}

	const { data: profile } = await supabaseClient.from('recruiter').select('*').maybeSingle();

	if (!profile && path !== '/(app)/recruiter/create') {
		throw redirect(303, '/recruiter/create');
	}

	return {
		profile: {
			...profile
		}
	};
};
