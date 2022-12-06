import type { Actions, FormData } from './$types';
import { invalid, redirect } from '@sveltejs/kit';
import { getSupabase } from '@supabase/auth-helpers-sveltekit';

// check url is valid
const isValidUrl = (str: string): boolean => {
	try {
		new URL(str);
		return true;
	} catch (error) {
		// If the URL is invalid, an error will be thrown
		return false;
	}
};

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
			throw redirect(303, '/login');
		}
		const { request } = event;
		const data = await request.formData();

		const email = session?.user?.email;
		const firstName = getTrimmedFormValue(data, 'firstName');
		const lastName = getTrimmedFormValue(data, 'lastName');
		const linkedin = getTrimmedFormValue(data, 'linkedin');
		const referrer = getTrimmedFormValue(data, 'referrer');
		const comment = getTrimmedFormValue(data, 'comment');

		// form validation
		// Trims whitespace for all fields
		if (!email) {
			return invalid(400, {
				errors: {
					email: 'Email is required'
				}
			});
		}

		if (!firstName) {
			return invalid(400, {
				errors: {
					firstName: 'First name is required'
				}
			});
		}
		if (!lastName) {
			return invalid(400, {
				errors: {
					lastName: 'Last name is required'
				}
			});
		}

		// require linkedin && valid URL && url contains linkedin
		if (!linkedin) {
			return invalid(400, {
				errors: {
					linkedin: 'LinkedIn Profile is required'
				}
			});
		} else if (!isValidUrl(linkedin)) {
			return invalid(400, {
				errors: {
					linkedin: 'LinkedIn Profile must be a valid URL'
				}
			});
		} else if (!linkedin.includes('linkedin')) {
			return invalid(400, {
				errors: {
					linkedin: 'LinkedIn Profile must be a valid LinkedIn Profile URL'
				}
			});
		}

		// require referrer
		if (!referrer) {
			return invalid(400, {
				errors: {
					referrer: 'Referrer is required'
				}
			});
		}
		
		//
		// add the user to waitlist
		//
		const survey = {
			referrer,
			comment
		};

		// insert into waitlist table

		return { success: true };
	}
};
