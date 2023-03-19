import type { SupabaseClient } from '@supabase/supabase-js';

import { refreshAccessToken } from '$lib/server/google/oauth';

// get a refreshed google access token
// tbd where the best location for this function is
export const getRefreshedGoogleAccessToken = async (client: SupabaseClient): Promise<string> => {
	// get google refresh token
	// relies on RLS to only return the refresh token for the current user
	const { data } = await client
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
