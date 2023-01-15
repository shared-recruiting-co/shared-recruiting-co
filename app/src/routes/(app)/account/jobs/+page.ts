import type { PageLoad } from './$types';
import { getSupabase } from '@supabase/auth-helpers-sveltekit';
import { redirect, error } from '@sveltejs/kit';

type Job = {
	job_id: string;
	company: string;
	job_title: string;
	emailed_at: string;
	recruiter: string;
	recruiter_email: string;
};

type Pagination = {
	page: number;
	perPage: number;
	numPages: number;
	numResults: number;
	hasNext: boolean;
	hasPrev: boolean;
};

type Data = {
	jobs: Job[];
	pagination: Pagination;
};

const PAGE_SIZE = 10;

export const load: PageLoad<Data> = async (event) => {
	const { session, supabaseClient } = await getSupabase(event);

	// require user to be logged in
	if (!session) {
		throw redirect(303, '/login');
	}

	// get the query parameters from the URL (default to 1)
	const page = parseInt(event.url.searchParams.get('page') || '1') || 1;
	const start = (page - 1) * PAGE_SIZE;
	const stop = start + PAGE_SIZE;

	// get jobs from database
	const { data: jobs, error: jobsError } = await supabaseClient
		.from('user_email_job')
		.select('*')
		.order('emailed_at', { ascending: false })
		.range(start, stop);

	// TODO: decide on error handling
	if (jobsError) {
		console.error(jobsError);
		throw error(500, jobsError.message);
	}

	const { count, error: countError } = await supabaseClient.from('user_email_job').select('*', {
		head: true,
		count: 'exact'
	});

	if (countError || !count) {
		console.error(countError);
		throw error(500, countError.message);
	}

	const formattedJobs = jobs.map(({ data, ...job }) => ({
		...job,
		emailed_at: new Date(job.emailed_at).toLocaleDateString(),
		// TS is complaining because jsonb can technically be string, number, boolean, or JSON, or JSON[]
		recruiter: data?.recruiter,
		recruiter_email: data?.recruiter_email
	}));

	return {
		jobs: formattedJobs,
		pagination: {
			page,
			perPage: PAGE_SIZE,
			numPages: Math.ceil(count / PAGE_SIZE),
			numResults: count,
			hasPrev: page > 1,
			hasNext: page < Math.ceil(count / PAGE_SIZE)
		}
	};
};
