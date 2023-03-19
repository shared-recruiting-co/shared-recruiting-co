import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

import { stop } from '$lib/server/google/gmail';
import { getRefreshedGoogleAccessToken } from '$lib/supabase/client.server';

export const POST: RequestHandler = async ({ locals: { getSession, supabase } }) => {
	const session = await getSession();
	if (!session) throw error(401, 'unauthorized');

	// get google refresh token
	let accessToken = '';
	try {
		accessToken = await getRefreshedGoogleAccessToken(supabase);
	} catch (err) {
		// do something
		if (err instanceof Error) {
			throw error(500, err);
		}

		throw error(500, 'unexpected error occurred');
	}

	// stop watching for new emails
	const stopResponse = await stop(accessToken);
	if (stopResponse.status !== 200) throw error(500, 'failed to unsubscribe to gmail notifications');

	// "deactivate" the email in db
	const { error: updateError } = await supabase
		.from('user_profile')
		.update({ is_active: false })
		.eq('user_id', session?.user.id);

	// should we re-subscribe if the db update fails?
	if (updateError) throw error(500, 'failed to saved changes to database');

	return new Response('success');
};
