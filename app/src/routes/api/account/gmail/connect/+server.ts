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

// TODO: Redirect with error message if in redirect flow
const connectEmail = async ({
	method,
	code,
	session,
	supabaseClient,
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
			console.error("invalid response from google");
		}
		// log error message for debugging and return user facing error
		throw error(400, "Failed to connect email. Please try again.");
	}

	// get json 
	const { access_token, refresh_token, scope, expires_in, token_type, id_token } = await tokenResponse.json();

	// validate scope
	if (!scope.includes(GMAIL_MODIFY_SCOPE)) {
		// log error message for debugging and return user facing error
		console.log(`user did not grant ${GMAIL_MODIFY_SCOPE} scope`);
		throw error(400, "Missing required scope. Please check all boxes on the consent screen.");
	}

	// parse id token to get user id
	const idTokenPayload = parseJWTPayload(id_token);
	if (!idTokenPayload) {
		console.log("invalid id token");
		throw error(400, "Something went wrong. Please try again.");
	}

	// verify emails match
	const { email } = idTokenPayload;
	if (email !== session.user.email) {
		console.log("authorized user does not match session user");
		throw error(400, `You must connect with the same email as your SRC account. Please connect with ${session.user.email}`);
	}

	// check if user has an account 
	const { data: profile } = await supabaseClient.from('user_profile').select('*').maybeSingle();
	// if no profile, create one
	if (!profile) {
	console.log('creating user profile...');
		// make sure the user is allowed to create one
		const { data: waitlist } = await supabaseClient.from('waitlist').select('*').maybeSingle();
		if (!waitlist) {
			console.log('user is not on waitlist');
			throw error(401, `${session.user.email} is not on the waitlist`);
		}
		if (!waitlist.can_create_account) {
			console.log( `user ${session.user.email} is not allowed to create an account yet`);
			throw error(401, `${session.user.email} is not allowed to create an account yet`);
		}
		// create a profile for them
		const { error: createError } = await adminSupabaseClient.from('user_profile').insert({
			user_id: session.user.id,
			first_name: waitlist.first_name,
			last_name: waitlist.last_name,
			email: session.user.email,
		})

		if (createError) {
			console.error("failed to create profile: ", createError);
			throw error(500, 'Failed to create your account. Please try again. If this problem persists, please team@sharedrecruiting.co.');
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
		console.error("failed to save oauth tokens:", tokenSaveError);
		throw error(500, 'Failed to sync your account. Please try again. If this problem persists, reach out to team@sharedrecruiting.co.');
	}

	// watch for new emails
	console.log('subscribing to gmail notifications...');
	const watchResponse = await watch(access_token);
	if (watchResponse.status !== 200) {
		console.log('failed to subscribe to gmail notifications');
		throw error(500, 'Failed to sync your account. Please try again. If this problem persists, reach out to team@sharedrecruiting.co.');
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
	
	await connectEmail({ method: "POST", code: code.toString(), session, supabaseClient });

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
	
	await connectEmail({ method: "GET", code, session, supabaseClient });

  return new Response("success")
};
