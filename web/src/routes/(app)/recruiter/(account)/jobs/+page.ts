import type { PageLoad } from './$types';
import { redirect, error } from '@sveltejs/kit';

import type { Database } from '$lib/supabase/types';

import { getPagePagination } from '$lib/pagination';
import type { Pagination } from '$lib/pagination';

type Job = Database['public']['Tables']['job']['Row'] & {
	candidate_count: number;
};

type Data = {
	jobs: Job[];
	pagination: Pagination;
};

export const load: PageLoad<Data> = async ({ url, parent }) => {
	const { session, supabase } = await parent();
	// require user to be logged in
	if (!session) {
		throw redirect(303, '/recruiter/login');
	}

	// get total jobs count from database
	const { count, error: countError } = await supabase
	.from('job')
	.select('*', {
		head: true,
		count: 'exact'
	});

	if (countError || count === null) {
		throw error(500, countError?.message);
	}


	// get pagination results 
	const pagination: Pagination = getPagePagination(url, count);

	// get jobs from database
	const { data: jobs, error: jobsError } = await supabase
		.from('job')
		.select('*')
		.order('updated_at', { ascending: false })
		.range(pagination.resultsToFetchStart, pagination.resultsToFetchEnd);

	// TODO: decide on error handling
	if (jobsError) {
		console.error(jobsError);
		throw error(500, jobsError.message);
	}

	// candidate counts
	const { data: candidateCounts, error: candidateCountsError } = await supabase
		.from('job_candidate_count')
		.select('*')
		.in(
			'job_id',
			jobs.map((job) => job.job_id)
		);

	// add candidate counts to jobs
	const jobsWithCandidateCounts = jobs.map((job) => {
		const candidateCount = candidateCounts.find(
			(candidateCount) => candidateCount.job_id === job.job_id
		)?.num_candidates;
		return {
			...job,
			candidate_count: candidateCount || 0
		};
	});

	return {jobs: jobsWithCandidateCounts, pagination};
};
