import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

import { PUBLIC_GOOGLE_CLIENT_ID, PUBLIC_GOOGLE_REDIRECT_URI } from '$env/static/public';
import { GOOGLE_CLIENT_SECRET } from '$env/static/private';
import { getSupabase } from '@supabase/auth-helpers-sveltekit';

import { supabaseClient as adminSupabaseClient } from '$lib/supabase/client.server';
import { exchangeCodeForTokens  } from '$lib/server/google/oauth';
import { watch  } from '$lib/server/google/gmail';

const GMAIL_MODIFY_SCOPE = 'https://www.googleapis.com/auth/gmail.modify';

// convert to ms and get iso string
const expriryFromExpiresIn = (expiresIn: number) => new Date(Date.now() + expiresIn * 1000).toISOString();

const parseJWTPayload = (token: string): Record<string,string> | null=> {
	try {
		return JSON.parse(Buffer.from(token.split('.')[1], 'base64').toString());
	} catch {
		return null
	}
}

 
export const POST: RequestHandler = async (event) => {
	// Create an callback endpoint that
	// - Confirm the X-Requested-With: XmlHttpRequest header is set for popup mode.
	// - Confirm CSRF token is valid.
	// - triggers a sync if first time
	// - subscribes to notifications
	// - validate on mobile!!
	// https://developers.google.com/identity/protocols/oauth2/web-server#httprest_3

	const { session, supabaseClient } = await getSupabase(event);
	if (!session) throw error(401, 'unauthorized');

  const { request } = event;
	const { headers } = request;

	console.log('headers', headers);

	const form = await request.formData();
	const code = form.get('code');
	if (!code) {
		throw error(400, 'missing code parameter');
	}
	
	// https://developers.google.com/identity/protocols/oauth2/web-server#httprest_3
	// exchange code for access and refresh token
  const tokenResponse = await exchangeCodeForTokens(code.toString());

	// what about redirects?
	if (tokenResponse.status !== 200) {
		throw error(400, 'failed to exchange code for access token');
	}

	// get json 
	const { access_token, refresh_token, scope, expires_in, token_type, id_token } = await tokenResponse.json();

	// validate scope
	if (!scope.includes(GMAIL_MODIFY_SCOPE)) {
		throw error(400, `user did not grant ${GMAIL_MODIFY_SCOPE} scope`);
	}

	// parse id token to get user id
	const idTokenPayload = parseJWTPayload(id_token);
	if (!idTokenPayload) {
		throw error(400, 'invalid id token');
	}

	// verify emails match
	const { email } = idTokenPayload;
	if (email !== session.user.email) {
		throw error(400, 'authorized user does not match session user');
	}

	// check if user has an account 
	const { data: profile } = await supabaseClient.from('user_profile').select('*').maybeSingle();
	// if no profile, create one
	if (!profile) {
	console.log('creating user profile...');
		// make sure the user is allowed to create one
		const { data: waitlist } = await supabaseClient.from('waitlist').select('*').maybeSingle();
		if (!waitlist) {
			throw error(401, `user ${session.user.email} is not on the waitlist`);
		}
		if (!waitlist.can_create_account) {
			throw error(401, `user ${session.user.email} is not allowed to create an account yet`);
		}
		// create a profile for them
		const { error: createError } = await adminSupabaseClient.from('user_profile').insert({
			user_id: session.user.id,
			first_name: waitlist.first_name,
			last_name: waitlist.last_name,
			email: session.user.email,
		})

		if (createError) {
			console.error(createError);
			throw error(500, 'failed to create user profile');
		}
	}

	// save oauth tokens
	console.log('saving user oauth token...');
	// format tokens for db
	const { error: tokenSaveError } = await supabaseClient.from('user_oauth_token').upsert({
		user_id: session.user.id,
		provider: 'google',
		token: {
			access_token,
			refresh_token,
			id_token,
			token_type,
			expiry: expriryFromExpiresIn(expires_in)
		}
	});

	if (tokenSaveError) {
		console.error(tokenSaveError);
		throw error(500, 'failed to save user oauth token');
	}

	// watch for new emails
	console.log('subscribing to gmail notifications...');
	const watchResponse = await watch(access_token);
	if (watchResponse.status !== 200) {
		throw error(500, 'failed to subscribe to gmail notifications');
	}

  return new Response("success")
};
