export type Json = string | number | boolean | null | { [key: string]: Json } | Json[];

export interface Database {
	graphql_public: {
		Tables: {
			[_ in never]: never;
		};
		Views: {
			[_ in never]: never;
		};
		Functions: {
			graphql: {
				Args: {
					operationName: string;
					query: string;
					variables: Json;
					extensions: Json;
				};
				Returns: Json;
			};
		};
		Enums: {
			[_ in never]: never;
		};
	};
	public: {
		Tables: {
			company: {
				Row: {
					company_name: string;
					website: string;
					company_id: string;
					created_at: string;
					updated_at: string;
				};
				Insert: {
					company_name: string;
					website: string;
					company_id?: string;
					created_at?: string;
					updated_at?: string;
				};
				Update: {
					company_name?: string;
					website?: string;
					company_id?: string;
					created_at?: string;
					updated_at?: string;
				};
			};
			job: {
				Row: {
					title: string;
					description_url: string;
					recruiter_id: string;
					company_id: string;
					job_id: string;
					created_at: string;
					updated_at: string;
				};
				Insert: {
					title: string;
					description_url: string;
					recruiter_id: string;
					company_id: string;
					job_id?: string;
					created_at?: string;
					updated_at?: string;
				};
				Update: {
					title?: string;
					description_url?: string;
					recruiter_id?: string;
					company_id?: string;
					job_id?: string;
					created_at?: string;
					updated_at?: string;
				};
			};
			recruiter: {
				Row: {
					user_id: string;
					email: string;
					first_name: string;
					last_name: string;
					company_id: string;
					responses: Json;
					created_at: string;
					updated_at: string;
				};
				Insert: {
					user_id: string;
					email: string;
					first_name: string;
					last_name: string;
					company_id: string;
					responses?: Json;
					created_at?: string;
					updated_at?: string;
				};
				Update: {
					user_id?: string;
					email?: string;
					first_name?: string;
					last_name?: string;
					company_id?: string;
					responses?: Json;
					created_at?: string;
					updated_at?: string;
				};
			};
			user_email_job: {
				Row: {
					user_id: string;
					user_email: string;
					email_thread_id: string;
					emailed_at: string;
					company: string;
					job_title: string;
					job_id: string;
					data: Json;
					created_at: string;
					updated_at: string;
				};
				Insert: {
					user_id: string;
					user_email: string;
					email_thread_id: string;
					emailed_at: string;
					company: string;
					job_title: string;
					job_id?: string;
					data?: Json;
					created_at?: string;
					updated_at?: string;
				};
				Update: {
					user_id?: string;
					user_email?: string;
					email_thread_id?: string;
					emailed_at?: string;
					company?: string;
					job_title?: string;
					job_id?: string;
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
					email: string | null;
					inbox_type: Database['public']['Enums']['inbox_type'];
					user_id: string;
					history_id: number;
					created_at: string;
					updated_at: string;
					synced_at: string;
				};
				Insert: {
					email?: string | null;
					inbox_type?: Database['public']['Enums']['inbox_type'];
					user_id: string;
					history_id: number;
					created_at?: string;
					updated_at?: string;
					synced_at?: string;
				};
				Update: {
					email?: string | null;
					inbox_type?: Database['public']['Enums']['inbox_type'];
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
					created_at: string;
					updated_at: string;
					token: Json;
					is_valid: boolean;
				};
				Insert: {
					user_id: string;
					provider: string;
					created_at?: string;
					updated_at?: string;
					token: Json;
					is_valid?: boolean;
				};
				Update: {
					user_id?: string;
					provider?: string;
					created_at?: string;
					updated_at?: string;
					token?: Json;
					is_valid?: boolean;
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
					auto_archive: boolean;
					auto_contribute: boolean;
					is_active: boolean;
				};
				Insert: {
					user_id: string;
					email: string;
					first_name: string;
					last_name: string;
					updated_at?: string;
					created_at?: string;
					auto_archive?: boolean;
					auto_contribute?: boolean;
					is_active?: boolean;
				};
				Update: {
					user_id?: string;
					email?: string;
					first_name?: string;
					last_name?: string;
					updated_at?: string;
					created_at?: string;
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
			candidate_oauth_token: {
				Row: {
					user_id: string | null;
					provider: string | null;
					token: Json | null;
					created_at: string | null;
					updated_at: string | null;
					is_valid: boolean | null;
				};
			};
			recruiter_oauth_token: {
				Row: {
					user_id: string | null;
					provider: string | null;
					token: Json | null;
					created_at: string | null;
					updated_at: string | null;
					is_valid: boolean | null;
				};
			};
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
			inbox_type: 'candidate' | 'recruiter';
		};
	};
	storage: {
		Tables: {
			buckets: {
				Row: {
					id: string;
					name: string;
					owner: string | null;
					created_at: string | null;
					updated_at: string | null;
					public: boolean | null;
				};
				Insert: {
					id: string;
					name: string;
					owner?: string | null;
					created_at?: string | null;
					updated_at?: string | null;
					public?: boolean | null;
				};
				Update: {
					id?: string;
					name?: string;
					owner?: string | null;
					created_at?: string | null;
					updated_at?: string | null;
					public?: boolean | null;
				};
			};
			migrations: {
				Row: {
					id: number;
					name: string;
					hash: string;
					executed_at: string | null;
				};
				Insert: {
					id: number;
					name: string;
					hash: string;
					executed_at?: string | null;
				};
				Update: {
					id?: number;
					name?: string;
					hash?: string;
					executed_at?: string | null;
				};
			};
			objects: {
				Row: {
					bucket_id: string | null;
					name: string | null;
					owner: string | null;
					metadata: Json | null;
					id: string;
					created_at: string | null;
					updated_at: string | null;
					last_accessed_at: string | null;
					path_tokens: string[] | null;
				};
				Insert: {
					bucket_id?: string | null;
					name?: string | null;
					owner?: string | null;
					metadata?: Json | null;
					id?: string;
					created_at?: string | null;
					updated_at?: string | null;
					last_accessed_at?: string | null;
					path_tokens?: string[] | null;
				};
				Update: {
					bucket_id?: string | null;
					name?: string | null;
					owner?: string | null;
					metadata?: Json | null;
					id?: string;
					created_at?: string | null;
					updated_at?: string | null;
					last_accessed_at?: string | null;
					path_tokens?: string[] | null;
				};
			};
		};
		Views: {
			[_ in never]: never;
		};
		Functions: {
			extension: {
				Args: { name: string };
				Returns: string;
			};
			filename: {
				Args: { name: string };
				Returns: string;
			};
			foldername: {
				Args: { name: string };
				Returns: string[];
			};
			get_size_by_bucket: {
				Args: Record<PropertyKey, never>;
				Returns: { size: number; bucket_id: string }[];
			};
			search: {
				Args: {
					prefix: string;
					bucketname: string;
					limits: number;
					levels: number;
					offsets: number;
					search: string;
					sortcolumn: string;
					sortorder: string;
				};
				Returns: {
					name: string;
					id: string;
					updated_at: string;
					created_at: string;
					last_accessed_at: string;
					metadata: Json;
				}[];
			};
		};
		Enums: {
			[_ in never]: never;
		};
	};
}
