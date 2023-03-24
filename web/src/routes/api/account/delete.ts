import type { SupabaseClient } from '@supabase/supabase-js';

import { refreshAccessToken } from '$lib/server/google/oauth';
import { sendMessage } from '$lib/server/google/gmail';

// naively sanitize user input
const sanitizeInput = (input: string) => {
	return input
		.replaceAll(/</g, '&lt;')
		.replaceAll(/>/g, '&gt;')
		.replaceAll(/script/gi, 'scr_ipt');
};

const founder = 'devin@sharedrecruiting.co';
const from = `Devin <${founder}>`;

// Send delete email with reason to founder
// Hack: Send email from founder account for now...
// TODO: Use a real transactional email service (sendgrid/mailgun) instead of this homebrew solution
export const sendDeleteEmail = async ({
	supabaseAdmin,
	userEmail,
	reason
}: {
	supabaseAdmin: SupabaseClient;
	userEmail: string;
	reason: string;
}) => {
	// get user id from pofile
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
		// send to self
		to: from,
		subject: 'SRC Account Deleted',
		body: `User: ${userEmail}\nReason: ${sanitizeInput(reason)}`
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
