import type { PageServerLoad } from './$types';
import { redirect, error } from '@sveltejs/kit';
import { getPagePagination } from '$lib/pagination';
import type { Pagination } from '$lib/pagination';



import type { JobInterest } from '$lib/supabase/client';

type Job = {
	job_id: string;
	job_title: string;
	job_description_url: string;
	job_interest: JobInterest
	company_name: string;
	company_website: string;
	recruiter_name: string;
	recruiter_email: string;
	emailed_at: string;
	is_verified: boolean;
};

type Data = {
	jobs: Job[];
	pagination: Pagination;
};

const PAGE_SIZE = 10;

export const load: PageServerLoad<Data> = async ({
	url,
	locals: { supabaseAdmin, getSession }
}) => {
	const session = await getSession();
	// require user to be logged in
	if (!session) {
		throw redirect(303, '/login');
	}

	// retreive the total count of the users jobs
	const { count, countError } = await	supabaseAdmin
		.from('vw_job_board')
		.select('*', {
			head: true,
			count: 'exact'
		})
		.eq('user_id', session.user.id)

	if (countError || count === null) {
		throw error(500, countError?.message);
	}

	// get pagination results 
	const pagination: Pagination = getPagePagination(url, count);

	// get jobs from database
	// HACK:
	// Use supabaseAdmin as a workaround for RLS not allowing candidates to see recruiter's name, email, and company info
	const { data: jobs, jobsError } = await supabaseAdmin
			.from('vw_job_board')
			.select('*')
			.eq('user_id', session.user.id)
			.order('emailed_at', { ascending: false })
			.range(pagination.resultsToFetchStart, pagination.resultsToFetchEnd)

	if (jobsError) {
		throw error(500, jobsError.message);
	}

	return {jobs, pagination};
};
