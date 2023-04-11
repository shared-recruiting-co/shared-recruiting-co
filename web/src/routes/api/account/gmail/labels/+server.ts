import { error, json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { getRefreshedGoogleAccessToken } from '$lib/supabase/client.server';
import { getSRCLabels, updateThreadLabels, isValidThread } from '$lib/server/google/gmail';
import type { Label } from '$lib/google/labels';

type DeleteRequest = {
	email: string;
	threadId: string;
};

// DELETE
// Remove all SRC gmail labels
export const DELETE: RequestHandler = async ({ request, locals: { getSession, supabase } }) => {
	// Get the current user's session
	const session = await getSession();
	if (!session) throw error(401, 'unauthorized');

	const { email, threadId } = (await request.json()) as DeleteRequest;
	if (!email) {
		throw error(400, 'email is required');
	}
	if (!threadId) {
		throw error(400, 'threadId is required');
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
	const srcLabels = await getSRCLabels(accessToken);
	const srcLabelIds = Object.values(srcLabels) as string[];

	// Remove all "@SRC" labels from the specified Gmail thread
	const success = await updateThreadLabels({ accessToken, threadId, removeLabelIds: srcLabelIds });

	return json({ success });
};

type PutRequest = {
	email: string;
	threadId: string;
	addLabels?: string[];
	removeLabels?: string[];
};

// PUT
// Update SRC gmail labels
export const PUT: RequestHandler = async ({ request, locals: { getSession, supabase } }) => {
	// Get the current user's session
	const session = await getSession();
	if (!session) throw error(401, 'unauthorized');

	const { email, threadId, addLabels, removeLabels } = (await request.json()) as PutRequest;
	if (!email) {
		throw error(400, 'email is required');
	}
	if (!threadId) {
		throw error(400, 'threadId is required');
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
	const srcLabels = await getSRCLabels(accessToken);
	// get label ids from srcLabels, filter out any undefined labels
	const addLabelIds = addLabels?.map((label) => srcLabels[label as Label]).filter(Boolean) || [];
	const removeLabelIds =
		removeLabels?.map((label) => srcLabels[label as Label]).filter(Boolean) || [];

	// Remove all "@SRC" labels from the specified Gmail thread
	const success = await updateThreadLabels({ accessToken, threadId, removeLabelIds, addLabelIds });

	return json({ success });
};
