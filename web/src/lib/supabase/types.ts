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
					operationName?: string;
					query?: string;
					variables?: Json;
					extensions?: Json;
				};
				Returns: Json;
			};
		};
		Enums: {
			[_ in never]: never;
		};
		CompositeTypes: {
			[_ in never]: never;
		};
	};
	public: {
		Tables: {
			company: {
				Row: {
					company_id: string;
					company_name: string;
					created_at: string;
					updated_at: string;
					website: string;
				};
				Insert: {
					company_id?: string;
					company_name: string;
					created_at?: string;
					updated_at?: string;
					website: string;
				};
				Update: {
					company_id?: string;
					company_name?: string;
					created_at?: string;
					updated_at?: string;
					website?: string;
				};
			};
			job: {
				Row: {
					company_id: string;
					created_at: string;
					description_url: string;
					job_id: string;
					recruiter_id: string;
					title: string;
					updated_at: string;
				};
				Insert: {
					company_id: string;
					created_at?: string;
					description_url: string;
					job_id?: string;
					recruiter_id: string;
					title: string;
					updated_at?: string;
				};
				Update: {
					company_id?: string;
					created_at?: string;
					description_url?: string;
					job_id?: string;
					recruiter_id?: string;
					title?: string;
					updated_at?: string;
				};
			};
			recruiter: {
				Row: {
					company_id: string;
					created_at: string;
					email: string;
					first_name: string;
					last_name: string;
					responses: Json;
					updated_at: string;
					user_id: string;
				};
				Insert: {
					company_id: string;
					created_at?: string;
					email: string;
					first_name: string;
					last_name: string;
					responses?: Json;
					updated_at?: string;
					user_id: string;
				};
				Update: {
					company_id?: string;
					created_at?: string;
					email?: string;
					first_name?: string;
					last_name?: string;
					responses?: Json;
					updated_at?: string;
					user_id?: string;
				};
			};
			user_email_job: {
				Row: {
					company: string;
					created_at: string;
					data: Json;
					email_thread_id: string;
					emailed_at: string;
					job_id: string;
					job_title: string;
					updated_at: string;
					user_email: string;
					user_id: string;
				};
				Insert: {
					company: string;
					created_at?: string;
					data?: Json;
					email_thread_id: string;
					emailed_at: string;
					job_id?: string;
					job_title: string;
					updated_at?: string;
					user_email: string;
					user_id: string;
				};
				Update: {
					company?: string;
					created_at?: string;
					data?: Json;
					email_thread_id?: string;
					emailed_at?: string;
					job_id?: string;
					job_title?: string;
					updated_at?: string;
					user_email?: string;
					user_id?: string;
				};
			};
			user_email_stat: {
				Row: {
					created_at: string;
					email: string;
					stat_id: string;
					stat_value: number;
					updated_at: string;
					user_id: string;
				};
				Insert: {
					created_at?: string;
					email: string;
					stat_id: string;
					stat_value?: number;
					updated_at?: string;
					user_id: string;
				};
				Update: {
					created_at?: string;
					email?: string;
					stat_id?: string;
					stat_value?: number;
					updated_at?: string;
					user_id?: string;
				};
			};
			user_email_sync_history: {
				Row: {
					created_at: string;
					email: string;
					history_id: number;
					inbox_type: Database['public']['Enums']['inbox_type'];
					synced_at: string;
					updated_at: string;
					user_id: string;
				};
				Insert: {
					created_at?: string;
					email: string;
					history_id: number;
					inbox_type: Database['public']['Enums']['inbox_type'];
					synced_at?: string;
					updated_at?: string;
					user_id: string;
				};
				Update: {
					created_at?: string;
					email?: string;
					history_id?: number;
					inbox_type?: Database['public']['Enums']['inbox_type'];
					synced_at?: string;
					updated_at?: string;
					user_id?: string;
				};
			};
			user_oauth_token: {
				Row: {
					created_at: string;
					email: string | null;
					is_valid: boolean;
					provider: string;
					token: Json;
					updated_at: string;
					user_id: string;
				};
				Insert: {
					created_at?: string;
					email?: string | null;
					is_valid?: boolean;
					provider: string;
					token: Json;
					updated_at?: string;
					user_id: string;
				};
				Update: {
					created_at?: string;
					email?: string | null;
					is_valid?: boolean;
					provider?: string;
					token?: Json;
					updated_at?: string;
					user_id?: string;
				};
			};
			user_profile: {
				Row: {
					auto_archive: boolean;
					auto_contribute: boolean;
					created_at: string;
					email: string;
					first_name: string;
					is_active: boolean;
					last_name: string;
					updated_at: string;
					user_id: string;
				};
				Insert: {
					auto_archive?: boolean;
					auto_contribute?: boolean;
					created_at?: string;
					email: string;
					first_name: string;
					is_active?: boolean;
					last_name: string;
					updated_at?: string;
					user_id: string;
				};
				Update: {
					auto_archive?: boolean;
					auto_contribute?: boolean;
					created_at?: string;
					email?: string;
					first_name?: string;
					is_active?: boolean;
					last_name?: string;
					updated_at?: string;
					user_id?: string;
				};
			};
			waitlist: {
				Row: {
					can_create_account: boolean;
					created_at: string;
					email: string;
					first_name: string;
					last_name: string;
					linkedin_url: string;
					responses: Json;
					updated_at: string;
					user_id: string;
				};
				Insert: {
					can_create_account?: boolean;
					created_at?: string;
					email: string;
					first_name: string;
					last_name: string;
					linkedin_url: string;
					responses?: Json;
					updated_at?: string;
					user_id: string;
				};
				Update: {
					can_create_account?: boolean;
					created_at?: string;
					email?: string;
					first_name?: string;
					last_name?: string;
					linkedin_url?: string;
					responses?: Json;
					updated_at?: string;
					user_id?: string;
				};
			};
		};
		Views: {
			candidate_oauth_token: {
				Row: {
					created_at: string | null;
					is_valid: boolean | null;
					provider: string | null;
					token: Json | null;
					updated_at: string | null;
					user_id: string | null;
				};
			};
			recruiter_oauth_token: {
				Row: {
					created_at: string | null;
					is_valid: boolean | null;
					provider: string | null;
					token: Json | null;
					updated_at: string | null;
					user_id: string | null;
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
		CompositeTypes: {
			[_ in never]: never;
		};
	};
	storage: {
		Tables: {
			buckets: {
				Row: {
					created_at: string | null;
					id: string;
					name: string;
					owner: string | null;
					public: boolean | null;
					updated_at: string | null;
				};
				Insert: {
					created_at?: string | null;
					id: string;
					name: string;
					owner?: string | null;
					public?: boolean | null;
					updated_at?: string | null;
				};
				Update: {
					created_at?: string | null;
					id?: string;
					name?: string;
					owner?: string | null;
					public?: boolean | null;
					updated_at?: string | null;
				};
			};
			migrations: {
				Row: {
					executed_at: string | null;
					hash: string;
					id: number;
					name: string;
				};
				Insert: {
					executed_at?: string | null;
					hash: string;
					id: number;
					name: string;
				};
				Update: {
					executed_at?: string | null;
					hash?: string;
					id?: number;
					name?: string;
				};
			};
			objects: {
				Row: {
					bucket_id: string | null;
					created_at: string | null;
					id: string;
					last_accessed_at: string | null;
					metadata: Json | null;
					name: string | null;
					owner: string | null;
					path_tokens: string[] | null;
					updated_at: string | null;
				};
				Insert: {
					bucket_id?: string | null;
					created_at?: string | null;
					id?: string;
					last_accessed_at?: string | null;
					metadata?: Json | null;
					name?: string | null;
					owner?: string | null;
					path_tokens?: string[] | null;
					updated_at?: string | null;
				};
				Update: {
					bucket_id?: string | null;
					created_at?: string | null;
					id?: string;
					last_accessed_at?: string | null;
					metadata?: Json | null;
					name?: string | null;
					owner?: string | null;
					path_tokens?: string[] | null;
					updated_at?: string | null;
				};
			};
		};
		Views: {
			[_ in never]: never;
		};
		Functions: {
			extension: {
				Args: {
					name: string;
				};
				Returns: string;
			};
			filename: {
				Args: {
					name: string;
				};
				Returns: string;
			};
			foldername: {
				Args: {
					name: string;
				};
				Returns: string[];
			};
			get_size_by_bucket: {
				Args: Record<PropertyKey, never>;
				Returns: {
					size: number;
					bucket_id: string;
				}[];
			};
			search: {
				Args: {
					prefix: string;
					bucketname: string;
					limits?: number;
					levels?: number;
					offsets?: number;
					search?: string;
					sortcolumn?: string;
					sortorder?: string;
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
		CompositeTypes: {
			[_ in never]: never;
		};
	};
}
