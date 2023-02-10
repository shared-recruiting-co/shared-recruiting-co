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

	const { data: profile } = await supabaseClient
		.from('recruiter')
		.select('*,company(company_id,company_name,website)')
		.maybeSingle();

	if (!profile) {
		if (path !== '/(app)/recruiter/create') {
			throw redirect(303, '/recruiter/create');
		}
		return {};
	}

	// if were aren't already on an account page, redirect to the profile page
	if (!path.includes('(account)/')) {
		throw redirect(303, '/recruiter/profile');
	}

	return {
		profile: {
			email: profile.email,
			firstName: profile.first_name,
			lastName: profile.last_name,
			createdAt: profile.create_at,
			updatedAt: profile.updated_at
		},
		company: {
			id: profile.company.company_id,
			name: profile.company.company_name,
			website: profile.company.website
		}
	};
};