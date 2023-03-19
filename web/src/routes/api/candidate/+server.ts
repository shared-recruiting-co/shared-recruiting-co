import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

import { sendWelcomeEmail } from '$lib/server/google/welcomeEmail';

export const POST: RequestHandler = async ({
	request,
	locals: { supabase, getSession, supabaseAdmin }
}) => {
	const session = await getSession();
	if (!session) throw error(401, 'unauthorized');

	const body = await request.json();
	const { tos } = body;

	// verify gmail is connected
	const { error: oauthError } = await supabase.from('user_oauth_token').select('*').maybeSingle();

	if (oauthError) {
		console.error('failed to get oauth token:', error);
		throw error(401, `Email must be connected to create an account`);
	}

	// verify tos is set to true
	if (!tos) {
		console.log('user did not agree to terms of service');
		throw error(400, 'You must agree to the terms of service to create an account.');
	}
	console.log('creating user profile...');
	// make sure the user is allowed to create one
	const { data: waitlist } = await supabase.from('waitlist').select('*').maybeSingle();
	if (!waitlist) {
		console.log('user is not on waitlist');
		throw error(401, `${session.user.email} is not on the waitlist`);
	}
	if (!waitlist.can_create_account) {
		console.log(`user ${session.user.email} is not allowed to create an account yet`);
		throw error(401, `${session.user.email} is not allowed to create an account yet`);
	}

	// create a profile for them
	const { error: createError } = await supabaseAdmin.from('user_profile').insert({
		user_id: session.user.id,
		first_name: waitlist.first_name,
		last_name: waitlist.last_name,
		email: session.user.email,
		// for now, all early users opt-in to auto contributing emails
		auto_contribute: true
	});

	if (createError) {
		console.error('failed to create profile: ', createError);
		throw error(
			500,
			'Failed to create your account. Please try again. If this problem persists, please team@sharedrecruiting.co.'
		);
	}

	const resp = await event.fetch('/api/account/gmail/subscribe', {
		method: 'POST'
	});

	// still create account
	if (!resp.ok) {
		console.log('failed to subscribe to gmail notifications');
		return new Response('success');
	}

	// send message from founder email to welcome (back) user
	// to trigger the initial email sync
	console.log('sending welcome email...');

	// TODO: Use a real transactional email service (sendgrid/mailgun) instead of this homebrew solution
	await sendWelcomeEmail({
		supabaseAdmin,
		email: session.user.email,
		isNewUser: true
	});

	return new Response('success');
};
