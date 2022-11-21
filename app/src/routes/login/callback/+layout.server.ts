import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

import { getSupabase } from '@supabase/auth-helpers-sveltekit';

// convert to ms and get iso string
const expriryFromExpiresAt = (expiresAt: number) => new Date(expiresAt * 1000).toISOString();

export const load: LayoutServerLoad = async (event) => {
	const { session, supabaseClient } = await getSupabase(event);
	if (!session) return {};

	// synchronize provider tokens with db
	const { expires_at, provider_token, provider_refresh_token, user } = session;
	const { provider } = user.app_metadata;

	if (!expires_at || !provider_token || !provider_refresh_token || !provider) {
		console.log('missing data. not updating user oauth token');
	} else {
		console.log('saving user oauth token...');
		// format tokens for db
		const { error } = await supabaseClient.from('user_oauth_token').upsert({
			user_id: user.id,
			provider,
			token: {
				access_token: provider_token,
				refresh_token: provider_refresh_token,
				expiry: expriryFromExpiresAt(expires_at)
			}
		});
		if (error) {
			console.log('failed to save user oauth token:', error);
		} else {
			// check if user has had their email synced
			const { data, error: selectError } = await supabaseClient
				.from('user_email_sync_history')
				.select('user_id')
				.maybeSingle();

			// if not, trigger login workflow
			if (!data && !selectError) {
				console.log('user has not had their email synced. triggering new user workflow...');
				// make an authenticated request to the login workflow
				// TODO
			}
		}
	}
	// redirect home
	throw redirect(307, '/');
};
