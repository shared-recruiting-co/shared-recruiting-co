import type { SupabaseClient } from '@supabase/supabase-js';

import { refreshAccessToken } from '$lib/server/google/oauth';
import { sendMessage } from '$lib/server/google/gmail';

const founder = 'devin@sharedrecruiting.co';
const from = `Devin <${founder}>`;

const newUserBody = `Hi 👋,

Welcome to the Shared Recruiting Co. This email triggers an one-time historic sync of your Gmail account. You should see your new @SRC labels in your Gmail account shortly. Moving forward, every time you receive an email about a job opportunity it will automatically be labeled and managed by SRC.

If you have any question or run into issues, just reply to this email. I'm here to help.

-Devin`;

const returningUserBody = `Hi 👋,

Welcome back to Shared Recruiting Co. We're happy to have you back. This email triggers a sync of your emails between now and the last time your account was active.

If you have any question or run into issues, just reply to this email. I'm here to help.

-Devin`;

// Hack: Send email from founder account for now...
// TODO: Use a real transactional email service (sendgrid/mailgun) instead of this homebrew solution
export const sendWelcomeEmail = async ({
	supabaseAdmin,
	email,
	isNewUser
}: {
	supabaseAdmin: SupabaseClient;
	email: string;
	isNewUser: boolean;
}) => {
	// get user id from profile
	const { data: profileData, error: profileError } = await supabaseAdmin
		.from('user_profile')
		.select('user_id')
		.eq('email', founder)
		.maybeSingle();

	if (!profileData || !profileData.user_id) {
		console.error('failed to get founder profile:', profileError);
		// fail silently for now
		return;
	}

	const { data } = await supabaseAdmin
		.from('user_oauth_token')
		.select('token')
		.eq('provider', 'google')
		.eq('user_id', profileData.user_id)
		.maybeSingle();

	if (!data || !data.token) {
		console.error('failed to get founder oauth token');
		// fail silently for now
		return;
	}

	const { token } = data;

	if (!token.refresh_token) {
		console.error('failed to get founder refresh token');
		// fail silently for now
		return;
	}

	// refresh google access token
	const resp = await refreshAccessToken(token.refresh_token);

	if (resp.status !== 200) {
		const error = await resp.json();
		console.error('failed to refresh founder access token:', error);
		// fail silently for now
		return;
	}

	const { access_token } = await resp.json();

	// send email
	const msg = {
		from,
		to: email,
		subject: isNewUser ? 'Welcome to the Shared Recruiting Co.' : 'Welcome Back to SRC',
		body: isNewUser ? newUserBody : returningUserBody
	};

	const sendResponse = await sendMessage(access_token, msg);
	if (sendResponse.status !== 200) {
		const error = await sendResponse.json();
		console.error('failed to send welcome email:', error);
		// fail silently for now
		return;
	}

	return;
};
