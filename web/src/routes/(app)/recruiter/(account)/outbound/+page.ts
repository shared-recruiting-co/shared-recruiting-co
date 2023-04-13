import type { PageLoad } from './$types';

import type { Database } from '$lib/supabase/types';

import { getPagePagination } from '$lib/pagination';
import type { Pagination } from '$lib/pagination';

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

	// get recruiter outbound template count for pagination
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

	// get pagination results 
	const pagination: Pagination = getPagePagination(url, count);

	// Get outbound recruiter template
	const { data: outboundTemplates } = await supabase
	.from('recruiter_outbound_template')
	.select('*,job(*)')
	.order('created_at', { ascending: false })
	.range(pagination.resultsToFetchStart, pagination.resultsToFetchEnd);

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
		pagination
	};
};
