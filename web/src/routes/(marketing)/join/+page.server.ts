import type { Actions } from './$types';
import { fail, redirect } from '@sveltejs/kit';

import { isValidUrl, getTrimmedFormValue } from '$lib/forms';

export const actions: Actions = {
	default: async (event) => {
		const {
			request,
			locals: { getSession, supabaseAdmin }
		} = event;
		const session = await getSession();
		// require user to be logged in
		if (!session) {
			throw redirect(303, '/login');
		}
		const data = await request.formData();

		const userId = session?.user?.id;
		const email = session?.user?.email;
		const firstName = getTrimmedFormValue(data, 'firstName');
		const lastName = getTrimmedFormValue(data, 'lastName');
		const linkedin = getTrimmedFormValue(data, 'linkedin');
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

		// require linkedin && valid URL && url contains linkedin
		if (!linkedin) {
			return fail(400, {
				errors: {
					linkedin: 'LinkedIn Profile is required'
				}
			});
		} else if (!isValidUrl(linkedin)) {
			return fail(400, {
				errors: {
					linkedin: 'LinkedIn Profile must be a valid URL'
				}
			});
		} else if (!linkedin.includes('linkedin')) {
			return fail(400, {
				errors: {
					linkedin: 'LinkedIn Profile must be a valid LinkedIn Profile URL'
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
			linkedin_url: linkedin,
			responses,
			can_create_account: false
		};

		const { error } = await supabaseAdmin.from('waitlist').insert(row);

		if (error) {
			console.error('error adding user to the waitlist', error);
			return fail(400, {
				errors: {
					submit: error.message
				}
			});
		}
		return { success: true };
	}
};
