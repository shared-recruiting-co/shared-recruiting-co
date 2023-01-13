export type Json = string | number | boolean | null | { [key: string]: Json } | Json[];

export interface Database {
	public: {
		Tables: {
			user_email_job: {
				Row: {
					job_id: string;
					user_id: string;
					user_email: string;
					email_thread_id: string;
					emailed_at: string;
					company: string;
					job_title: string;
					data: Json;
					created_at: string;
					updated_at: string;
				};
				Insert: {
					job_id?: string;
					user_id: string;
					user_email: string;
					email_thread_id: string;
					emailed_at: string;
					company: string;
					job_title: string;
					data?: Json;
					created_at?: string;
					updated_at?: string;
				};
				Update: {
					job_id?: string;
					user_id?: string;
					user_email?: string;
					email_thread_id?: string;
					emailed_at?: string;
					company?: string;
					job_title?: string;
					data?: Json;
					created_at?: string;
					updated_at?: string;
				};
			};
			user_email_stat: {
				Row: {
					user_id: string;
					email: string;
					stat_id: string;
					stat_value: number;
					created_at: string;
					updated_at: string;
				};
				Insert: {
					user_id: string;
					email: string;
					stat_id: string;
					stat_value?: number;
					created_at?: string;
					updated_at?: string;
				};
				Update: {
					user_id?: string;
					email?: string;
					stat_id?: string;
					stat_value?: number;
					created_at?: string;
					updated_at?: string;
				};
			};
			user_email_sync_history: {
				Row: {
					user_id: string;
					history_id: number;
					created_at: string;
					updated_at: string;
					synced_at: string;
				};
				Insert: {
					user_id: string;
					history_id: number;
					created_at?: string;
					updated_at?: string;
					synced_at?: string;
				};
				Update: {
					user_id?: string;
					history_id?: number;
					created_at?: string;
					updated_at?: string;
					synced_at?: string;
				};
			};
			user_oauth_token: {
				Row: {
					user_id: string;
					provider: string;
					token: Json;
					created_at: string;
					updated_at: string;
					is_valid: boolean;
				};
				Insert: {
					user_id: string;
					provider: string;
					token: Json;
					created_at?: string;
					updated_at?: string;
					is_valid?: boolean;
				};
				Update: {
					user_id?: string;
					provider?: string;
					token?: Json;
					created_at?: string;
					updated_at?: string;
					is_valid?: boolean;
				};
			};
			user_profile: {
				Row: {
					user_id: string;
					email: string;
					first_name: string;
					last_name: string;
					created_at: string;
					updated_at: string;
					auto_archive: boolean;
					auto_contribute: boolean;
					is_active: boolean;
				};
				Insert: {
					user_id: string;
					email: string;
					first_name: string;
					last_name: string;
					created_at?: string;
					updated_at?: string;
					auto_archive?: boolean;
					auto_contribute?: boolean;
					is_active?: boolean;
				};
				Update: {
					user_id?: string;
					email?: string;
					first_name?: string;
					last_name?: string;
					created_at?: string;
					updated_at?: string;
					auto_archive?: boolean;
					auto_contribute?: boolean;
					is_active?: boolean;
				};
			};
			waitlist: {
				Row: {
					user_id: string;
					email: string;
					first_name: string;
					last_name: string;
					linkedin_url: string;
					responses: Json;
					can_create_account: boolean;
					created_at: string;
					updated_at: string;
				};
				Insert: {
					user_id: string;
					email: string;
					first_name: string;
					last_name: string;
					linkedin_url: string;
					responses?: Json;
					can_create_account?: boolean;
					created_at?: string;
					updated_at?: string;
				};
				Update: {
					user_id?: string;
					email?: string;
					first_name?: string;
					last_name?: string;
					linkedin_url?: string;
					responses?: Json;
					can_create_account?: boolean;
					created_at?: string;
					updated_at?: string;
				};
			};
		};
		Views: {
			[_ in never]: never;
		};
		Functions: {
			increment_user_email_stat: {
				Args: {
					user_id: string;
					email: string;
					stat_id: string;
					stat_value: number;
				};
				Returns: undefined;
			};
		};
		Enums: {
			[_ in never]: never;
		};
	};
}
