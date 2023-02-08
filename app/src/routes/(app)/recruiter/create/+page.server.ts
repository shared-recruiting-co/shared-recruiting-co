import type { Actions, FormData } from './$types';
import { fail, redirect } from '@sveltejs/kit';
import { getSupabase } from '@supabase/auth-helpers-sveltekit';

import { supabaseClient } from '$lib/supabase/client.server';

const getTrimmedFormValue = (data: FormData, key: string): string => {
	const value = data.get(key);
	if (typeof value === 'string') {
		return value.trim();
	}
	return '';
};

export const actions: Actions = {
	default: async (event) => {
		const { session } = await getSupabase(event);
		// require user to be logged in
		if (!session) {
			throw redirect(303, '/recruiter/login');
		}
		const { request } = event;
		const data = await request.formData();

		const userId = session?.user?.id;
		const email = session?.user?.email;
		const firstName = getTrimmedFormValue(data, 'firstName');
		const lastName = getTrimmedFormValue(data, 'lastName');
		const company = getTrimmedFormValue(data, 'company');
		const referrer = getTrimmedFormValue(data, 'referrer');
		const comment = getTrimmedFormValue(data, 'comment');

		// form validation
		// Trims whitespace for all fields
		if (!email) {
			return fail(400, {
				errors: {
					email: 'Email is required'
				}
			});
		}

		if (!firstName) {
			return fail(400, {
				errors: {
					firstName: 'First name is required'
				}
			});
		}
		if (!lastName) {
			return fail(400, {
				errors: {
					lastName: 'Last name is required'
				}
			});
		}

		// require company
		if (!company) {
			return fail(400, {
				errors: {
					company: 'Company is required'
				}
			});
		}

		// require referrer
		if (!referrer) {
			return fail(400, {
				errors: {
					referrer: 'Referrer is required'
				}
			});
		}

		const responses = {
			referrer,
			comment
		};

		// add the user to waitlist
		const row = {
			user_id: userId,
			email,
			first_name: firstName,
			last_name: lastName,
			company_name: company,
			responses
		};

		const { error } = await supabaseClient.from('recruiter').insert(row);

		if (error) {
			console.error('error creating recruiter', error);
			return fail(400, {
				errors: {
					submit: error.message
				}
			});
		}

		return { success: true };
	}
};
