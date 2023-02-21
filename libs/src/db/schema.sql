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

-- Enable Suapbase Realtime
begin;
  -- remove the supabase_realtime publication
  drop publication if exists supabase_realtime;

  -- re-create the supabase_realtime publication with no tables
  create publication supabase_realtime;
commit;

-- User OAuth Token Table
create table public.user_oauth_token (
    user_id uuid references auth.users(id) on delete cascade not null,
    email text not null,
    provider text not null,
    token jsonb not null,
    is_valid boolean not null default true,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),

    primary key (user_id, email, provider)
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

create type inbox_type as enum ('candidate', 'recruiter');

-- User Email Sync History Table
create table public.user_email_sync_history (
    user_id uuid references auth.users(id) on delete cascade not null,
    inbox_type inbox_type not null,
    email text not null,
    history_id int8 not null,
    -- track successful sync attempts
    synced_at timestamp with time zone not null default now(),
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),

    primary key (user_id, inbox_type, email)
);

create trigger handle_updated_at_user_email_sync_history before update on public.user_email_sync_history
  for each row execute procedure moddatetime (updated_at);

-- lock down permissions on public.user_email_sync_history
-- only service role can insert. user cannot modify
alter table public.user_email_sync_history enable row level security;

create policy "Users can view their own email sync history"
  on public.user_email_sync_history for select
  using ( auth.uid() = user_id );

-- enable realtime
alter publication supabase_realtime add table user_email_sync_history;

-- Waitlist table
create table public.waitlist (
    user_id uuid references auth.users(id) on delete cascade not null,
    -- duplicate email for convenience
    email text not null,
    first_name text not null,
    last_name text not null,
    linkedin_url text not null,
    responses jsonb not null default '{}'::jsonb,
    can_create_account boolean not null default false,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),

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
    user_id uuid references auth.users(id) on delete cascade not null,
    -- duplicate email for convenience
    email text not null,
    first_name text not null,
    last_name text not null,
    -- user email settings
    -- in future, this will be migrated to a separate table to support multi-email accounts
    is_active boolean not null default true,
    auto_archive boolean not null default false,
    auto_contribute boolean not null default false,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),

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

-- user_email_stat table
-- simple table to keep track of realtime user facing statistics
create table public.user_email_stat (
    user_id uuid references auth.users(id) on delete cascade not null,
    -- in future, user can have multiple emails
    email text not null,
    -- free form stat id
    stat_id text not null,
    -- integer for now
    stat_value integer not null default 0,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),

    primary key (user_id, email, stat_id)
);

create trigger handle_updated_at_user_email_stat before update on public.user_email_stat
  for each row execute procedure moddatetime (updated_at);

-- enable RLS
alter table public.user_email_stat enable row level security;
create policy "Users can view their own email stats"
  on public.user_email_stat for select
  using ( auth.uid() = user_id );

-- enable realtime
alter publication supabase_realtime add table user_email_stat;

create or replace function increment_user_email_stat (user_id uuid, email text, stat_id text, stat_value int)
returns void as
$$
  insert into public.user_email_stat(user_id, email, stat_id, stat_value)
  values (user_id, email, stat_id, stat_value)
  on conflict (user_id, email, stat_id)
  do update set
      stat_value = user_email_stat.stat_value + excluded.stat_value;
$$
language sql volatile;


-- user_email_job
create table public.user_email_job (
  job_id uuid not null default uuid_generate_v4(),
  user_id uuid references auth.users(id) on delete cascade not null,
  user_email text not null,
  -- thread id of the recruiting email
  email_thread_id text not null,
  -- time when the email was delivered to user
  emailed_at timestamp with time zone not null,
  company text not null,
  job_title text not null,
  -- use catch-all jsonb for everything else while we figure out the job schema
  data jsonb not null default '{}'::jsonb,
  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone not null default now(),

  -- prevent duplicate email_thread_id processing
  unique (user_email, email_thread_id),
  -- prevent duplicate job posting
  -- for now consider, (user_id, company, job_title) to be a unique job posting
  unique (user_id, company, job_title),
  primary key (job_id)
);

create trigger handle_updated_at_user_email_job before update on public.user_email_job
  for each row execute procedure moddatetime (updated_at);

-- enable realtime
alter publication supabase_realtime add table user_email_job;

-- enable RLS
alter table public.user_email_job enable row level security;
create policy "Users can view their own jobs"
  on public.user_email_job for select
  using ( auth.uid() = user_id );

-- for now jobs are read only

-- company table
create table public.company (
    company_id uuid not null default uuid_generate_v4(),
    company_name text not null,
    website text not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),

    primary key (company_id)
);

create trigger handle_updated_at_company before update on public.company
  for each row execute procedure moddatetime (updated_at);

-- enable realtime
alter publication supabase_realtime add table public.company;
alter table public.company enable row level security;

-- recruiter table
-- similar to user_profile but for recruiters
create table public.recruiter (
    user_id uuid references auth.users(id) on delete cascade not null,
    -- duplicate email for convenience
    email text not null,
    first_name text not null,
    last_name text not null,
    -- collect additional information when recruiters create an account
    responses jsonb not null default '{}'::jsonb,

    company_id uuid references public.company(company_id) not null,

    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),

    primary key (user_id)
);

create trigger handle_updated_at_recruiter before update on public.recruiter
  for each row execute procedure moddatetime (updated_at);

-- enable realtime
alter publication supabase_realtime add table public.recruiter;
alter table public.recruiter enable row level security;

create policy "Recruiters can view their own profile"
  on public.recruiter for select
  using ( auth.uid() = user_id );

create policy "Recruiters can update their own profile"
  on public.recruiter for update
  using ( auth.uid() = user_id );

create policy "Recruiters can view their own company"
  on public.company for select
  using ( auth.uid() in (
      select user_id
      from public.company
      inner join recruiter using (company_id)
      where user_id = auth.uid()
  ));

create policy "Recruiters can view their own company"
  on public.company for select
  using ( company_id in (
      select company_id
      from public.recruiter
      where user_id = auth.uid()
  ));

-- jobs
create table public.job (
    job_id uuid not null default uuid_generate_v4(),
    title text not null,
    description_url text not null,

    recruiter_id uuid references public.recruiter(user_id) on delete cascade not null,
    company_id uuid references public.company(company_id) on delete cascade not null,

    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),

    primary key (job_id)
);

create trigger handle_updated_at_job before update on public.job
  for each row execute procedure moddatetime (updated_at);

-- enable real-time
alter publication supabase_realtime add table public.job;
alter table public.job enable row level security;

create policy "Recruiters can view their jobs"
  on public.job for select
  using ( auth.uid() = recruiter_id );

create policy "Recruiters can insert their jobs"
  on public.job for insert
  with check ( auth.uid() = recruiter_id );

create policy "Recruiters can update their jobs"
  on public.job for update
  using ( auth.uid() = recruiter_id );

create policy "Recruiters can delete their jobs"
  on public.job for delete
  using ( auth.uid() = recruiter_id );

create policy "Companies can view their jobs"
  on public.job for select
  using ( company_id in (
      select company_id 
      from public.recruiter
      where user_id = auth.uid()
  ));


-- views on oauth tokens
create view candidate_oauth_token as
select
  user_oauth_token.*
from user_oauth_token
inner join user_profile using (user_id);

create view recruiter_oauth_token as
select
  user_oauth_token.*
from user_oauth_token
inner join recruiter using (user_id);
