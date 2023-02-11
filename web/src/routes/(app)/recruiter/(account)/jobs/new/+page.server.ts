import type { Actions } from './$types';
import { fail, redirect } from '@sveltejs/kit';
import { getSupabase } from '@supabase/auth-helpers-sveltekit';

import { isValidUrl, getTrimmedFormValue } from '$lib/forms';

export const actions: Actions = {
	default: async (event) => {
		const { session, supabaseClient } = await getSupabase(event);
		// require user to be logged in
		if (!session) {
			throw redirect(303, '/recruiter/login');
		}
		const { request } = event;
		const data = await request.formData();

		const userId = session?.user?.id;
		const jobTitle = getTrimmedFormValue(data, 'title');
		const jobDescriptionURL = getTrimmedFormValue(data, 'descriptionURL');

		// form validation
		// Trims whitespace for all fields
		if (!jobTitle) {
			return fail(400, {
				errors: {
					title: 'Job title is required'
				}
			});
		}

		if (!jobDescriptionURL) {
			return fail(400, {
				errors: {
					descriptionURL: 'Job description URL is required'
				}
			});
		} else if (!isValidUrl(jobDescriptionURL)) {
			return fail(400, {
				errors: {
					descriptionURL: 'Job description must be a valid URL'
				}
			});
		}

		const { data: profile, error: profileError } = await supabaseClient
			.from('recruiter')
			.select('*')
			.maybeSingle();

		if (!profile) {
			console.error('no recruiter profile found: error', profileError);
		}

		// add job
		const row = {
			recruiter_id: userId,
			company_id: profile.company_id,
			title: jobTitle,
			description_url: jobDescriptionURL
		};

		const { error } = await supabaseClient.from('job').insert(row);

		if (error) {
			console.error('error creating job', error);
			return fail(400, {
				errors: {
					submit: error.message
				}
			});
		}

		return { success: true };
	}
};
