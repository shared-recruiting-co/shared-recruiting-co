import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

import { getSupabase } from '@supabase/auth-helpers-sveltekit';

import { supabaseClient as adminSupabaseClient } from '$lib/supabase/client.server';
import { exchangeCodeForTokens  } from '$lib/server/google/oauth';
import { watch  } from '$lib/server/google/gmail';
import type {Session, SupabaseClient} from '@supabase/supabase-js';

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

const connectEmail = async ({
	code,
	session,
	supabaseClient,
}: {
 code: string;
 session: Session;
 supabaseClient: SupabaseClient;
}) => {
	// https://developers.google.com/identity/protocols/oauth2/web-server#httprest_3
	// exchange code for access and refresh token
  const tokenResponse = await exchangeCodeForTokens(code);

	// what about redirects?
	if (tokenResponse.status !== 200) {
		const err = await tokenResponse.json();
		console.error(err);
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
}

export const POST: RequestHandler = async (event) => {
	const { session, supabaseClient } = await getSupabase(event);
	if (!session) throw error(401, 'unauthorized');

  const { request } = event;
	const { headers } = request;

	const  xRequestedWith = headers.get('x-requested-with') || '';
	if (xRequestedWith !== 'XmlHttpRequest') throw error(400, 'invalid request');

	const form = await request.formData();
	const code = form.get('code');
	if (!code) {
		throw error(400, 'missing code parameter');
	}
	
	await connectEmail({ code: code.toString(), session, supabaseClient });

  return new Response("success")
};

export const GET: RequestHandler = async (event) => {
	const { session, supabaseClient } = await getSupabase(event);
	if (!session) throw error(401, 'unauthorized');

  const {  url } = event;

	const code = url.searchParams.get('code');
	if (!code) {
		throw error(400, 'missing code parameter');
	}
	
	await connectEmail({ code, session, supabaseClient });

  return new Response("success")
};
