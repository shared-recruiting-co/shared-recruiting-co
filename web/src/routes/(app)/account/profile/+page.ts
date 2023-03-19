import type { PageLoad } from './$types';
import { redirect } from '@sveltejs/kit';
import { UserEmailStats } from '$lib/supabase/client';

type Data = {
	lastSyncedAt?: string;
	isSetup: boolean;
	numEmailsProcessed: number | null;
	numJobsDetected: number | null;
};

export const load: PageLoad<Data> = async ({ parent }) => {
	const { session, supabase } = await parent();
	// require user to be logged in
	if (!session) {
		throw redirect(303, '/login');
	}

	const [{ data: emailSyncHistory }, { data: oauthToken }, { data: emailStats }] =
		await Promise.all([
			supabase.from('user_email_sync_history').select('synced_at').maybeSingle(),
			supabase.from('user_oauth_token').select('is_valid').maybeSingle(),
			supabase.from('user_email_stat').select('*')
		]);

	//TODO: Move this aggregation to the database
	const numEmailsProcessed =
		emailStats?.reduce(
			(acc, stat) =>
				stat.stat_id === UserEmailStats.EmailsProcessed ? acc + stat.stat_value : acc,
			0
		) || 0;

	const numJobsDetected =
		emailStats?.reduce(
			(acc, stat) => (stat.stat_id === UserEmailStats.JobsDetected ? acc + stat.stat_value : acc),
			0
		) || 0;

	return {
		lastSyncedAt: emailSyncHistory?.synced_at as string | undefined,
		isSetup: Boolean(oauthToken?.is_valid),
		numEmailsProcessed,
		numJobsDetected
	};
};
