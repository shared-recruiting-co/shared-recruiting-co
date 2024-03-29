import type { PageLoad } from './$types';
import { redirect } from '@sveltejs/kit';

type Data = {
	lastSyncedAt?: string;
	isSetup: boolean;
	numInboundJobs: number | null;
	numOfficialJobs: number | null;
};

export const load: PageLoad<Data> = async ({ parent }) => {
	const { session, supabase } = await parent();
	// require user to be logged in
	if (!session) {
		throw redirect(303, '/login');
	}

	const [
		{ data: emailSyncHistory },
		{ data: oauthToken },
		{ data: candidateJobCountUnverified },
		{ data: candidateJobCount }
	] = await Promise.all([
		// get the last time the user synced any of their emails
		supabase
			.from('user_email_sync_history')
			.select('synced_at')
			.order('synced_at', { ascending: false })
			.limit(1)
			.maybeSingle(),
		supabase.from('user_oauth_token').select('is_valid').maybeSingle(),
		supabase.from('candidate_job_count_unverified').select('*').maybeSingle(),
		supabase.from('candidate_job_count').select('*').maybeSingle()
	]);

	return {
		lastSyncedAt: emailSyncHistory?.synced_at || (new Date().toISOString() as string | undefined),
		isSetup: Boolean(oauthToken?.is_valid),
		numInboundJobs: candidateJobCountUnverified?.num_jobs || 0,
		numOfficialJobs: candidateJobCount?.num_jobs || 0
	};
};
