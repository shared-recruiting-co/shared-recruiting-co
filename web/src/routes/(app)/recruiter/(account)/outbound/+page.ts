import type { PageLoad } from './$types';

import type { Database } from '$lib/supabase/types';

const PAGE_SIZE = 10;

type Pagination = {
	page: number;
	perPage: number;
	numPages: number;
	numResults: number;
	hasNext: boolean;
	hasPrev: boolean;
};

type Sequence = Database['public']['Tables']['recruiter_outbound_template']['Row'] & {
	recipient_count: number;
	job: Database['public']['Tables']['job']['Row'];
};

type Data = {
	pagination: Pagination;
	sequences: Sequence[];
	jobs: {
		job_id: string;
		title: string;
	}[];
};

export const load: PageLoad<Data> = async ({ url, parent }) => {
	const { supabase } = await parent();

	// assume parent is checking for session
	const page = parseInt(url.searchParams.get('page') || '1') || 1;
	const start = (page - 1) * PAGE_SIZE;

	const stop = start + PAGE_SIZE;

	// TODO: Use Promise.all() to get all the data at once

	// Get outbound recruiter template
	const { data: outboundTemplates } = await supabase
		.from('recruiter_outbound_template')
		.select('*,job(*)')
		.order('created_at', { ascending: false })
		.range(start, stop);

	let { count, error: countError } = await supabase
		.from('recruiter_outbound_template')
		.select('*', {
			head: true,
			count: 'exact'
		});

	// TODO: decide on error handling
	if (countError) {
		console.error(countError);
	}
	if (count === null) {
		count = 0;
	}

	// Get count of recipients (candidates) per template
	const { data: recipientCounts, error: recipientCountsError } = await supabase
		.from('outbound_template_recipient_count')
		.select('*')
		.in('template_id', outboundTemplates?.map((t) => t.template_id) || []);

	// update sequences with recipient counts
	const sequences = outboundTemplates?.map((sequence) => {
		const recipientCount = recipientCounts?.find(
			(recipientCount) => recipientCount.template_id === sequence.template_id
		)?.num_recipients;
		return {
			...sequence,
			recipient_count: recipientCount || 0
		};
	});

	// fetch possible jobs
	// TODO: Consider doing this client-side in the combo box
	const { data: jobs } = await supabase.from('job').select('job_id,title');

	return {
		sequences,
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
