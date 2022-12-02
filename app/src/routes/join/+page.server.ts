import type { Actions } from './$types';
import { invalid } from '@sveltejs/kit';
import { getSupabase } from '@supabase/auth-helpers-sveltekit';

export const actions: Actions = {
	default: async (event) => {
		const { session } = await getSupabase(event);
		const { request } = event;
		const data = await request.formData();
		const email = data.get('email');
		const firstName = data.get('firstName');
		const lastName = data.get('lastName');
		const linkedin = data.get('linkedin');
		const referrer = data.get('referrer');
		const comment = data.get('comment');
		// form validation
		if (!email || !firstName || !lastName) {
			return invalid(400, 'All fields are required');
		}
		// require firstName
		// require lastName
		// require linkedin && valid URL && url contains linkedin
		if (!linkedin) {
			return invalid(400, 'All fields are required');
		}
		// require referrer
		//
		// add the user to waitlist
		//
		const survey = {
			referrer,
			comment
		};

		// insert into waitlist table
	}
};
