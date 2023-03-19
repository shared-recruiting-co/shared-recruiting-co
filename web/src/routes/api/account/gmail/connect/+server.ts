import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

import { exchangeCodeForTokens } from '$lib/server/google/oauth';
import type { Session, SupabaseClient } from '@supabase/supabase-js';

const GMAIL_MODIFY_SCOPE = 'https://www.googleapis.com/auth/gmail.modify';

// convert to ms and get iso string
const expriryFromExpiresIn = (expiresIn: number) =>
	new Date(Date.now() + expiresIn * 1000).toISOString();

const parseJWTPayload = (token: string): Record<string, string> | null => {
	try {
		return JSON.parse(Buffer.from(token.split('.')[1], 'base64').toString());
	} catch {
		return null;
	}
};

// TODO: Redirect with error message if in redirect flow
const connectEmail = async ({
	method,
	code,
	session,
	supabaseClient
}: {
	method: 'GET' | 'POST';
	code: string;
	session: Session;
	supabaseClient: SupabaseClient;
}) => {
	// https://developers.google.com/identity/protocols/oauth2/web-server#httprest_3
	// exchange code for access and refresh token
	const tokenResponse = await exchangeCodeForTokens(code, method);

	// what about redirects?
	if (tokenResponse.status !== 200) {
		try {
			const err = await tokenResponse.json();
			console.error(err);
		} catch {
			console.error('invalid response from google');
		}
		// log error message for debugging and return user facing error
		throw error(400, 'Failed to connect email. Please try again.');
	}

	// get json
	const { access_token, refresh_token, scope, expires_in, token_type, id_token } =
		await tokenResponse.json();

	// validate scope
	if (!scope.includes(GMAIL_MODIFY_SCOPE)) {
		// log error message for debugging and return user facing error
		console.log(`user did not grant ${GMAIL_MODIFY_SCOPE} scope`);
		throw error(400, 'Missing required scope. Please check all boxes on the consent screen.');
	}

	// parse id token to get user id
	const idTokenPayload = parseJWTPayload(id_token);
	if (!idTokenPayload) {
		console.log('invalid id token');
		throw error(400, 'Something went wrong. Please try again.');
	}

	// verify emails match
	const { email } = idTokenPayload;

	// save oauth tokens
	console.log('saving user oauth token...');
	// format tokens for db
	const { error: tokenSaveError } = await supabaseClient.from('user_oauth_token').upsert({
		user_id: session.user.id,
		email,
		provider: 'google',
		token: {
			access_token,
			refresh_token,
			id_token,
			token_type,
			expiry: expriryFromExpiresIn(expires_in)
		},
		is_valid: true
	});

	if (tokenSaveError) {
		console.error('failed to save oauth tokens:', tokenSaveError);
		throw error(
			500,
			'Failed to sync your account. Please try again. If this problem persists, reach out to team@sharedrecruiting.co.'
		);
	}
};

export const POST: RequestHandler = async ({ request, locals: { getSession, supabase } }) => {
	const session = await getSession();
	if (!session) throw error(401, 'unauthorized');

	const { headers } = request;

	const xRequestedWith = headers.get('x-requested-with') || '';
	if (xRequestedWith !== 'XmlHttpRequest') throw error(400, 'invalid request');

	const form = await request.formData();
	const code = form.get('code');
	if (!code) {
		throw error(400, 'missing code parameter');
	}

	await connectEmail({ method: 'POST', code: code.toString(), session, supabase });

	return new Response('success');
};

export const GET: RequestHandler = async ({ url, locals: { getSession, supabase } }) => {
	const session = await getSession();
	if (!session) throw error(401, 'unauthorized');

	const code = url.searchParams.get('code');
	if (!code) {
		throw error(400, 'missing code parameter');
	}

	await connectEmail({ method: 'GET', code, session, supabase });

	return new Response('success');
};
