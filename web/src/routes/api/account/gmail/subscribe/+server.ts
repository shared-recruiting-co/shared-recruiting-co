import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

import { CANDIDATE_GMAIL_PUBSUB_TOPIC, RECRUITER_GMAIL_PUBSUB_TOPIC } from '$env/static/private';

import { watch } from '$lib/server/google/gmail';
import type { WatchRequest } from '$lib/server/google/gmail';
import { getRefreshedGoogleAccessToken } from '$lib/supabase/client.server';

const candidateWatchRequest: WatchRequest = {
	labelIds: ['UNREAD'],
	labelFilterAction: 'include',
	topicName: CANDIDATE_GMAIL_PUBSUB_TOPIC
};

const recruiterWatchRequest: WatchRequest = {
	labelIds: ['SENT'],
	labelFilterAction: 'include',
	topicName: RECRUITER_GMAIL_PUBSUB_TOPIC
};

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

	// check if user is a candidate or recruiter
	const { data: candidate } = await supabase.from('user_profile').select('user_id').maybeSingle();
	if (candidate) {
		// watch for new emails
		const watchResponse = await watch(accessToken, candidateWatchRequest);

		if (watchResponse.status !== 200)
			throw error(500, 'Failed to subscribe to gmail notifications');

		// "activate" the email in db
		const { error: updateError } = await supabase
			.from('user_profile')
			.update({ is_active: true })
			.eq('user_id', session?.user.id);

		if (updateError) throw error(500, 'failed to saved changes to database');
	}

	const { data: recruiter } = await supabase.from('recruiter').select('user_id').maybeSingle();
	if (recruiter) {
		// watch for new emails
		const watchResponse = await watch(accessToken, recruiterWatchRequest);

		if (watchResponse.status !== 200)
			throw error(500, 'Failed to subscribe to gmail notifications');

		// TODO
		// "activate" the email in db
		// const { error: updateError } = await supabase
		// .from('user_profile')
		// .update({ is_active: true })
		// .eq('user_id', session?.user.id);
		//
		// if (updateError) throw error(500, 'failed to saved changes to database');
	}

	return new Response('success');
};

// options
// check if request if from a candidate or recruiter (or both) -> subscribe to topic accordingly
//
// separate endpoint for updating email settings (/candidate/email_settings, /recruiter/email_settings)
// -> in this case, check if is_active is toggled to true from false, if so, send email
// create separate endpoint
