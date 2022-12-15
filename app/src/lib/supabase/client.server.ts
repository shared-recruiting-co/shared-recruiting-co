import { PUBLIC_SUPABASE_URL } from '$env/static/public';
import { SUPABASE_SERVICE_ROLE_KEY } from '$env/static/private';

import { createClient } from '@supabase/auth-helpers-sveltekit';
import type { SupabaseClient } from '@supabase/supabase-js';

import { refreshAccessToken } from '$lib/server/google/oauth';

// https://github.com/supabase/auth-helpers/tree/main/packages/sveltekit
export const supabaseClient = createClient(PUBLIC_SUPABASE_URL, SUPABASE_SERVICE_ROLE_KEY);

// get a refreshed google access token
// tbd where the best location for this function is
export const getRefreshedGoogleAccessToken = async (
	supabaseClient: SupabaseClient
): Promise<string> => {
	// get google refresh token
	// relies on RLS to only return the refresh token for the current user
	const { data } = await supabaseClient
		.from('user_oauth_token')
		.select('token')
		.eq('provider', 'google')
		.maybeSingle();

	if (!data || !data.token) throw new Error('No google oauth tokens found');

	const { token } = data;

	if (!token.refresh_token) throw new Error('No google refresh token found');

	// refresh google access token
	const resp = await refreshAccessToken(token.refresh_token);

	// TODO: Mark access token invalid if there is an ouauth error
	if (resp.status !== 200) {
		const error = await resp.json();
		throw new Error(
			`Failed to refresh google access token: ${error.error_description || error.message}`
		);
	}

	const { access_token } = await resp.json();

	return access_token;
};
