import type { PageLoad } from './$types';
import { error } from '@sveltejs/kit';

import type { Database } from '$lib/supabase/types';

type Sequence = Database['public']['Tables']['recruiter_outbound_template']['Row'] & {
	recipient_count: number;
};

type Data = {
	job: Database['public']['Tables']['job']['Row'];
	candidates: Database['public']['Tables']['candidate_company_inbound']['Row'][];
	outboundTemplates: Sequence[];
};

export const load: PageLoad<Data> = async ({ params, parent }) => {
	const { id: jobID } = params;
	const { supabase } = await parent();

	const { data: job, error: jobErr } = await supabase
		.from('job')
		.select('*')
		.eq('job_id', jobID)
		.single();

	if (!job) {
		throw error(404, 'Job not found');
	}
	if (jobErr) throw error(500, jobErr);

	// TODO:
	// Pagination!
	// Recipient Counts
	let [
		// eslint-disable-next-line
		{ data: candidates = [], error: candidatesErr },
		// eslint-disable-next-line
		{ data: outboundTemplates, error: outboundTemplateErr }
	] = await Promise.all([
		supabase.from('candidate_company_inbound').select('*').eq('job_id', jobID),
		supabase.from('recruiter_outbound_template').select('*').eq('job_id', jobID)
	]);

	if (candidatesErr) throw error(500, candidatesErr);
	if (!candidates) candidates = [];

	if (outboundTemplateErr) throw error(500, outboundTemplateErr);
	if (!outboundTemplates) outboundTemplates = [];

	const { data: recipientCounts, error: recipientCountsError } = await supabase
		.from('outbound_template_recipient_count')
		.select('*')
		.in('template_id', outboundTemplates?.map((t) => t.template_id) || []);

	if (recipientCountsError) {
		// for now just log
		console.error('error fetching outbound_template_recipient_count:', recipientCountsError);
	}

	const sequences = outboundTemplates?.map((sequence) => {
		const recipientCount = recipientCounts?.find(
			(recipientCount) => recipientCount.template_id === sequence.template_id
		)?.num_recipients;
		return {
			...sequence,
			recipient_count: recipientCount || 0
		};
	});

	return {
		job,
		candidates,
		outboundTemplates: sequences
	};
};
