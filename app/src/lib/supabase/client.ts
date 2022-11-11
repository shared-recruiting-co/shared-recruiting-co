import { PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY } from '$env/static/public';
import { createClient } from '@supabase/auth-helpers-sveltekit';

import { Database } from '$lib/supabase/types';

// https://github.com/supabase/auth-helpers/tree/main/packages/sveltekit
export const supabaseClient = createClient<Database>(PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY);
