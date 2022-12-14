import { PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY } from '$env/static/public';
import { createClient } from '@supabase/auth-helpers-sveltekit';

// https://github.com/supabase/auth-helpers/tree/main/packages/sveltekit
export const supabaseClient = createClient(PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY);

export enum UserEmailStats {
	EmailsProcessed = 'emails_processed',
	JobsDetected = 'jobs_detected'
}
