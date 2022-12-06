export type Json =
  | string
  | number
  | boolean
  | null
  | { [key: string]: Json }
  | Json[]

export interface Database {
  public: {
    Tables: {
      user_email_sync_history: {
        Row: {
          user_id: string
          history_id: number
          created_at: string | null
          updated_at: string | null
          examples_collected_at: string | null
        }
        Insert: {
          user_id: string
          history_id: number
          created_at?: string | null
          updated_at?: string | null
          examples_collected_at?: string | null
        }
        Update: {
          user_id?: string
          history_id?: number
          created_at?: string | null
          updated_at?: string | null
          examples_collected_at?: string | null
        }
      }
      user_oauth_token: {
        Row: {
          user_id: string
          provider: string
          token: Json | null
          created_at: string | null
          updated_at: string | null
          is_valid: boolean
        }
        Insert: {
          user_id: string
          provider: string
          token?: Json | null
          created_at?: string | null
          updated_at?: string | null
          is_valid?: boolean
        }
        Update: {
          user_id?: string
          provider?: string
          token?: Json | null
          created_at?: string | null
          updated_at?: string | null
          is_valid?: boolean
        }
      }
      user_profile: {
        Row: {
          user_id: string
          email: string
          first_name: string
          last_name: string
          created_at: string | null
          updated_at: string | null
        }
        Insert: {
          user_id: string
          email: string
          first_name: string
          last_name: string
          created_at?: string | null
          updated_at?: string | null
        }
        Update: {
          user_id?: string
          email?: string
          first_name?: string
          last_name?: string
          created_at?: string | null
          updated_at?: string | null
        }
      }
      waitlist: {
        Row: {
          user_id: string
          email: string
          first_name: string
          last_name: string
          linkedin_url: string
          responses: Json
          can_create_account: boolean
          created_at: string | null
          updated_at: string | null
        }
        Insert: {
          user_id: string
          email: string
          first_name: string
          last_name: string
          linkedin_url: string
          responses?: Json
          can_create_account?: boolean
          created_at?: string | null
          updated_at?: string | null
        }
        Update: {
          user_id?: string
          email?: string
          first_name?: string
          last_name?: string
          linkedin_url?: string
          responses?: Json
          can_create_account?: boolean
          created_at?: string | null
          updated_at?: string | null
        }
      }
    }
    Views: {
      [_ in never]: never
    }
    Functions: {
      [_ in never]: never
    }
    Enums: {
      [_ in never]: never
    }
  }
}

