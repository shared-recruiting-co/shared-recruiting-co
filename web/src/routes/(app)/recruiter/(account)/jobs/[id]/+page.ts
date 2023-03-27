import type { PageLoad } from './$types';
import { error } from '@sveltejs/kit';

import type { Database } from '$lib/supabase/types';

type Data = {
	job: Database['public']['Tables']['job']['Row'];
	candidates: Database['public']['Tables']['candidate_company_inbound']['Row'][];
	outboundTemplates: Database['public']['Tables']['recruiter_outbound_template']['Row'][];
};

export const load: PageLoad<Data> = async ({ params, parent }) => {
	const { id: jobID } = params;
	const { supabase } = await parent();

	// TODO: Use Promise.all() to get all the data at once
	const { data: job, error: jobErr } = await supabase
		.from('job')
		.select('*')
		.eq('job_id', jobID)
		.single();

	if (!job) {
		throw error(404, 'Job not found');
	}
	if (jobErr) throw error(500, jobErr);

	// TODO: Pagination!
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

	return {
		job,
		candidates,
		outboundTemplates
	};
};
