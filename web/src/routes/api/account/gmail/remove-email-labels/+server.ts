import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { getRefreshedGoogleAccessToken } from '$lib/supabase/client.server';
import { getSrcLabelIds, removeLabelsFromThread, isValidThread } from '$lib/server/google/gmail';

/**
 * Main handler function for this POST request route
 *
 * @param request - The incoming request object
 * @param locals.getSession - function used to retrieve the current user's session.
 * @param locals.supabase - the supabase client
 * @returns {Response} - A response indicating whether or not the operation was successful.
 * @throws Throws an error if the user is not authorized or if an unexpected error occurs.
 */
export const POST: RequestHandler = async ({ request, locals: { getSession, supabase } }) => {
	// Get the current user's session
	const session = await getSession();
	if (!session) throw error(401, 'unauthorized');

	// retrieve the necessary information from the request
	const { email, email_thread_id: threadId } = await request.json();

	if (!email || !threadId) {
		throw error(400, 'Could not find both email and email_thread_id from request json');
	}

	// Get a refreshed Google access token for the current user
	let accessToken = '';
	try {
		accessToken = await getRefreshedGoogleAccessToken(supabase, email);
	} catch (err) {
		throw error(500, err instanceof Error ? err : 'Could not retrieve access token');
	}

	// check the thread actually exists in Gmail before attempting the update
	const validThread = await isValidThread(accessToken, threadId);
	if (!validThread) {
		throw error(404, 'No corresponding Gmail thread');
	}

	// Get the label IDs for all "@SRC" labels
	const srcLabelIds = await getSrcLabelIds(accessToken);

	// Remove all "@SRC" labels from the specified Gmail thread
	const success = await removeLabelsFromThread(accessToken, threadId, srcLabelIds);

	// handle the return value of removeLabelsFromThread
	if (!success) {
		throw error(500, 'Failed to remove SRC labels');
	}

	// success
	return new Response('success');
};
