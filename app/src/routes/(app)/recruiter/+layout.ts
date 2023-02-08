import type { PageLoad } from './$types';
import { getSupabase } from '@supabase/auth-helpers-sveltekit';
import { redirect } from '@sveltejs/kit';

export const load: PageLoad = async (event) => {
	const { session, supabaseClient } = await getSupabase(event);

	// redirect to login
	if (!session) {
		throw redirect(303, '/recruiter/login');
	}

	const { data: profile } = await supabaseClient.from('recruiter').select('*').maybeSingle();

	// TODO: Do we need to check for current path to avoid infinite redirect
	const path = event.route.id;
	console.log('path', path);

	if (!profile && path !== '/(app)/recruiter/create') {
		throw redirect(303, '/recruiter/create');
	}

	return {
		profile: {
			...profile
		}
	};
};
