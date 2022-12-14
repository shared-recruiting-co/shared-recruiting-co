-- Start: Supabase auth.users table
-- Include a simplified Supabase auth.users table in the schema for SQLC compilation
create schema if not exists auth;
create table auth.users (
    id uuid primary key,
    email text not null
);
-- End: Supabase auth.users table

-- Add moddatetime extension
create extension if not exists moddatetime schema extensions;

-- User OAuth Token Table
create table public.user_oauth_token (
    user_id uuid references auth.users(id) not null,
    provider text not null,
    token jsonb,
    is_valid boolean not null default true,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now(),

    primary key (user_id, provider)
);

create trigger handle_updated_at_user_oauth_token before update on public.user_oauth_token
  for each row execute procedure moddatetime (updated_at);

alter table public.user_oauth_token enable row level security;


create policy "Users can view their own oauth tokens"
  on public.user_oauth_token for select
  using ( auth.uid() = user_id );

create policy "Users can insert their own oauth tokens."
  on public.user_oauth_token for insert
  with check ( auth.uid() = user_id );

create policy "Users can update own oauth tokens."
  on public.user_oauth_token for update
  using ( auth.uid() = user_id );

-- User Email Sync History Table
create table public.user_email_sync_history (
    user_id uuid references auth.users(id) not null,
    history_id int8 not null,
    -- track successful sync attempts
    synced_at timestamp with time zone default now(),
    -- temporary column to track when we've last collected examples for this user
    examples_collected_at date,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now(),

    primary key (user_id)
);

create trigger handle_updated_at_user_email_sync_history before update on public.user_email_sync_history
  for each row execute procedure moddatetime (updated_at);

-- lock down permissions on public.user_email_sync_history
-- only service role can insert. user cannot modify
alter table public.user_email_sync_history enable row level security;

create policy "Users can view their own email sync history"
  on public.user_email_sync_history for select
  using ( auth.uid() = user_id );

-- waitlist table
create table public.waitlist (
    user_id uuid references auth.users(id) not null,
    -- duplicate email for convenience
    email text not null,
    first_name text not null,
    last_name text not null,
    linkedin_url text not null,
    responses jsonb not null default '{}'::jsonb,
    can_create_account boolean not null default false,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now(),

    primary key (user_id)
);

create trigger handle_updated_at_waitlist before update on public.waitlist
  for each row execute procedure moddatetime (updated_at);

-- enable RLS so users can't modify they waitlist entry
alter table public.waitlist enable row level security;

create policy "Users can view their own waitlist entry"
  on public.waitlist for select
  using ( auth.uid() = user_id );

-- user_profile table
create table public.user_profile (
    user_id uuid references auth.users(id) not null,
    -- duplicate email for convenience
    email text not null,
    first_name text not null,
    last_name text not null,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now(),

    primary key (user_id)
);

create trigger handle_updated_at_profile before update on public.user_profile
  for each row execute procedure moddatetime (updated_at);

alter table public.user_profile enable row level security;

create policy "Users can view their own profile"
  on public.user_profile for select
  using ( auth.uid() = user_id );

create policy "Users can update their own profile"
  on public.user_profile for update
  using ( auth.uid() = user_id );
