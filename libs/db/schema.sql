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
