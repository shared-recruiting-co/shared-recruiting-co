/// <reference types="@sveltejs/kit" />
import type { SupabaseClient, Session } from '@supabase/supabase-js';
import type { Database } from './lib/supabase/types';

// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
// and what to do when importing types
declare namespace App {
	interface Supabase {
		Database: Database;
		SchemaName: 'public';
	}

	interface Locals {
		supabase: SupabaseClient<Database>;
		supabaseAdmin: SupabaseClient<Database>;
		getSession(): Promise<Session | null>;
	}
	interface PageData {
		supabase: SupabaseClient<Database>;
		session: Session | null;
	}
	interface Error {
		message: string;
		code: string;
	}
	// interface Platform {}
}
