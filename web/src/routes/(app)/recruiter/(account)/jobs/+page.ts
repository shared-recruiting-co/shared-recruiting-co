import type { PageLoad } from './$types';
import { redirect, error } from '@sveltejs/kit';

type Job = {
	job_id: string;
	company_id: string;
	title: string;
	description_url: string;
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

export const load: PageLoad<Data> = async ({ url, parent }) => {
	const { session, supabase } = await parent();
	// require user to be logged in
	if (!session) {
		throw redirect(303, '/recruiter/login');
	}

	// get the query parameters from the URL (default to 1)
	const page = parseInt(url.searchParams.get('page') || '1') || 1;
	const start = (page - 1) * PAGE_SIZE;
	const stop = start + PAGE_SIZE;

	// get jobs from database
	const { data: jobs, error: jobsError } = await supabase
		.from('job')
		.select('*')
		.order('updated_at', { ascending: false })
		.range(start, stop);

	// TODO: decide on error handling
	if (jobsError) {
		console.error(jobsError);
		throw error(500, jobsError.message);
	}

	const { count, error: countError } = await supabase.from('job').select('*', {
		head: true,
		count: 'exact'
	});

	if (countError || count === null) {
		throw error(500, countError?.message);
	}

	return {
		jobs,
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
