import { PUBLIC_SUPABASE_URL } from '$env/static/public';
import { SUPABASE_SERVICE_ROLE_KEY } from '$env/static/private';

import { createClient } from '@supabase/auth-helpers-sveltekit';

// https://github.com/supabase/auth-helpers/tree/main/packages/sveltekit
export const supabaseClient = createClient(PUBLIC_SUPABASE_URL, SUPABASE_SERVICE_ROLE_KEY);
