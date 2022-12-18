export type Json = string | number | boolean | null | { [key: string]: Json } | Json[];

export interface Database {
	public: {
		Tables: {
			user_email_sync_history: {
				Row: {
					user_id: string;
					history_id: number;
					synced_at: string;
					created_at: string;
					updated_at: string;
				};
				Insert: {
					user_id: string;
					history_id: number;
					synced_at?: string;
					created_at?: string;
					updated_at?: string;
				};
				Update: {
					user_id?: string;
					history_id?: number;
					synced_at?: string;
					created_at?: string;
					updated_at?: string;
				};
			};
			user_oauth_token: {
				Row: {
					user_id: string;
					provider: string;
					is_valid: boolean;
					created_at: string;
					updated_at: string;
					token: Json;
				};
				Insert: {
					user_id: string;
					provider: string;
					is_valid?: boolean;
					created_at?: string;
					updated_at?: string;
					token: Json;
				};
				Update: {
					user_id?: string;
					provider?: string;
					is_valid?: boolean;
					created_at?: string;
					updated_at?: string;
					token?: Json;
				};
			};
			user_profile: {
				Row: {
					user_id: string;
					email: string;
					first_name: string;
					last_name: string;
					updated_at: string;
					created_at: string;
					is_active: boolean;
					auto_archive: boolean;
					auto_contribute: boolean;
				};
				Insert: {
					user_id: string;
					email: string;
					first_name: string;
					last_name: string;
					updated_at?: string;
					created_at?: string;
					is_active?: boolean;
					auto_archive?: boolean;
					auto_contribute?: boolean;
				};
				Update: {
					user_id?: string;
					email?: string;
					first_name?: string;
					last_name?: string;
					updated_at?: string;
					created_at?: string;
					is_active?: boolean;
					auto_archive?: boolean;
					auto_contribute?: boolean;
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
			[_ in never]: never;
		};
		Enums: {
			[_ in never]: never;
		};
	};
}
