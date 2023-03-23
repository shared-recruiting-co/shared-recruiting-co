import type { PageLoad } from './$types';
import { error } from '@sveltejs/kit';

export const load: PageLoad = async ({ params, parent }) => {
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

	return {
		job
	};
};
