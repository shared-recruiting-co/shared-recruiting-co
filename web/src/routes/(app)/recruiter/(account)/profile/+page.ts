import type { PageLoad } from './$types';

import type { Database } from '$lib/supabase/types';

type Data = {
	gmailAccounts: Database['public']['Tables']['user_oauth_token']['Row'][];
};

export const load: PageLoad<Data> = async ({ parent }) => {
	// assume parent is checking for session
	const { supabase } = await parent();

	const { data: gmailAccounts } = await supabase
		.from('recruiter_oauth_token')
		.select('*')
		.eq('provider', 'google');

	if (!gmailAccounts) {
		return {
			gmailAccounts: []
		};
	}

	return {
		gmailAccounts
	};
};
