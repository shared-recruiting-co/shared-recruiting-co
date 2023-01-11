import type { PageLoad } from './$types';
import { getSupabase } from '@supabase/auth-helpers-sveltekit';
import { redirect } from '@sveltejs/kit';

type Job = {
	job_id: string;
	company: string;
	role: string;
	recruiter: string;
	recruiter_email: string;
	emailed_at: string;
};

type Data = {
	jobs: Job[];
};

export const load: PageLoad<Data> = async (event) => {
	const { session } = await getSupabase(event);

	// require user to be logged in
	if (!session) {
		throw redirect(303, '/login');
	}

	return {
		jobs: [
			{
				job_id: '1',
				company: 'Shared Recruiting Co.',
				role: 'Founding Engineer',
				recruiter: 'Devin Stein',
				recruiter_email: 'devin@sharedrecruiting.co',
				emailed_at: new Date().toLocaleDateString()
			}
		]
	};
};
