import type { RequestHandler } from './$types';
import { error, json } from '@sveltejs/kit';
import * as Sentry from '@sentry/node';

import { getSupabase } from '@supabase/auth-helpers-sveltekit';

import { stop } from '$lib/server/google/gmail';
import {
	supabaseClient as adminSupabaseClient,
	getRefreshedGoogleAccessToken
} from '$lib/supabase/client.server';

import { sendDeleteEmail } from './delete';

export const DELETE: RequestHandler = async (event) => {
	const { session, supabaseClient } = await getSupabase(event);
	if (!session) throw error(401, 'unauthorized');

	// validate request body has a 'reason'
	let { reason } = (await event.request.json()) || {};
	if (!reason || !reason.trim()) throw error(400, 'Reason is required');
	reason = reason.trim();

	// unsubscribe (if oauth is valid and not paused)
	// get google refresh token
	let accessToken = '';
	try {
		accessToken = await getRefreshedGoogleAccessToken(supabaseClient);
		// stop watching for new emails
		const stopResponse = await stop(accessToken);
		if (stopResponse.status !== 200)
			throw new Error(`error stopping gmail watch: status code ${stopResponse.status}`);
	} catch (err) {
		// log error but proceed to delete account
		console.log('error refreshing access token:', err);
		Sentry.captureException(err, { event });
	}
	// delete user
	const { error: deleteError } = await adminSupabaseClient.auth.admin.deleteUser(session.user.id);
	if (deleteError) {
		console.log('error deleting user:', deleteError);
		Sentry.captureException(deleteError, { event });
		throw error(
			500,
			'Failed to delete user. If this error persists, please contact team@sharedrecruiting.co'
		);
	}

	// sanitize and send email
	const { email } = session.user;
	try {
		await sendDeleteEmail({ userEmail: email, reason });
	} catch (err) {
		console.log('error sending delete email:', err);
		Sentry.captureException(err, { event });
	}

	return json({
		message: `user ${session.user.id} deleted`
	});
};
