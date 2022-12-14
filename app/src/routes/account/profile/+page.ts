import type { PageLoad } from './$types';
import { getSupabase } from '@supabase/auth-helpers-sveltekit';
import { redirect } from '@sveltejs/kit';

type Data = {
	lastSyncedAt?: string;
};

export const load: PageLoad<Data> = async (event) => {
	const { session, supabaseClient } = await getSupabase(event);

	// require user to be logged in
	if (!session) {
		throw redirect(303, '/login');
	}

	const { data: emailSyncHistory } = await supabaseClient
		.from('user_email_sync_history')
		.select('synced_at')
		.maybeSingle();

	// TODO: Check and return the status of the user oauth token

	return {
		lastSyncedAt: emailSyncHistory?.synced_at as string | undefined
	};
};
