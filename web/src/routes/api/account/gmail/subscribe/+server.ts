import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

import { watch } from '$lib/server/google/gmail';
import { getRefreshedGoogleAccessToken } from '$lib/supabase/client.server';

export const POST: RequestHandler = async ({ request, locals: { getSession, supabase } }) => {
	const session = await getSession();
	if (!session) throw error(401, 'unauthorized');

	// get email from request body
	let { email } = await request.json();
	email = email || session.user.email;

	// get google refresh token
	let accessToken = '';
	try {
		accessToken = await getRefreshedGoogleAccessToken(supabase, email);
	} catch (err) {
		// do something
		if (err instanceof Error) {
			throw error(500, err);
		}

		throw error(500, 'unexpected error occurred');
	}
	// watch for new emails
	const watchResponse = await watch(accessToken);

	if (watchResponse.status !== 200) throw error(500, 'Failed to subscribe to gmail notifications');

	// "activate" the email in db
	const { error: updateError } = await supabase
		.from('user_profile')
		.update({ is_active: true })
		.eq('user_id', session?.user.id);

	if (updateError) throw error(500, 'failed to saved changes to database');

	return new Response('success');
};

// options
// check if request if from a candidate or recruiter (or both) -> subscribe to topic accordingly
//
// separate endpoint for updating email settings (/candidate/email_settings, /recruiter/email_settings)
// -> in this case, check if is_active is toggled to true from false, if so, send email
// create separate endpoint
