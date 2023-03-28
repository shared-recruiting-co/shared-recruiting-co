import type { PageServerLoad } from './$types';
import { redirect, error } from '@sveltejs/kit';

type Job = {
	job_id: string;
	job_title: string;
	job_description_url: string;
	company_name: string;
	company_website: string;
	recruiter_name: string;
	recruiter_email: string;
	emailed_at: string;
	is_verified: boolean;
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

export const load: PageServerLoad<Data> = async ({
	url,
	locals: { supabaseAdmin, getSession }
}) => {
	const session = await getSession();
	// require user to be logged in
	if (!session) {
		throw redirect(303, '/login');
	}

	// get the query parameters from the URL (default to 1)
	const page = parseInt(url.searchParams.get('page') || '1') || 1;
	const start = (page - 1) * PAGE_SIZE;
	const stop = start + PAGE_SIZE;

	// get jobs from database
	// HACK:
	// Use supabaseAdmin as a workaround for RLS not allowing candidates to see recruiter's name, email, and company info
	const [{ data: jobs, jobsError }, { count, countError }] = await Promise.all([
		supabaseAdmin
			.from('vw_job_board')
			.select('*')
			.eq('user_id', session.user.id)
			.order('emailed_at', { ascending: false })
			.range(start, stop),
		supabaseAdmin
			.from('vw_job_board')
			.select('*', {
				head: true,
				count: 'exact'
			})
			.eq('user_id', session.user.id)
	]);
	if (jobsError) {
		throw error(500, jobsError.message);
	}
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
