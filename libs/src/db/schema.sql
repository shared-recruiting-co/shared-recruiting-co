--------------------------------
-- Start: Supabase auth.users table
--------------------------------
-- Include a simplified Supabase auth.users table in the schema for SQLC compilation
create schema if not exists auth;
create table auth.users (
    id uuid primary key,
    email text not null
);
--------------------------------
-- End: Supabase auth.users table
--------------------------------

--------------------------------
-- Start: Postgres Extensions
--------------------------------

-- Add moddatetime extension
create extension if not exists moddatetime schema extensions;

-- Add pg_trgm extension
create extension if not exists pg_trgm schema extensions;

-- Enable the pgtap extension for testing
create extension pgtap with schema extensions;

--------------------------------
-- End: Postgres Extensions
--------------------------------

-- Enable Suapbase Realtime
begin;
  -- remove the supabase_realtime publication
  drop publication if exists supabase_realtime;

  -- re-create the supabase_realtime publication with no tables
  create publication supabase_realtime;
commit;

--------------------------------
-- Start: User OAuth Token Table
--------------------------------
create table public.user_oauth_token (
    user_id uuid references auth.users(id) on delete cascade not null,
    email text not null,
    provider text not null,
    token jsonb not null,
    is_valid boolean not null default true,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),

    -- Add a unique constraint to prevent duplicate tokens
    unique (email, provider),
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
 
--------------------------------
-- End: User OAuth Token Table
--------------------------------

--------------------------------
-- Start: User Email Sync History Table
--------------------------------
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

--------------------------------
-- End: User Email Sync History Table
--------------------------------

--------------------------------
-- Start: Waitlist Table
--------------------------------

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

--------------------------------
-- End: Waitlist Table
--------------------------------

--------------------------------
-- Start: User (Candidate) Profile Table
--------------------------------

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

--------------------------------
-- End: User (Candidate) Profile Table
--------------------------------

--------------------------------
-- Start: User (Candidate) Email Stat Table
--------------------------------

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

--------------------------------
-- End: User (Candidate) Email Stat Table
--------------------------------

--------------------------------
-- Start: User (Candidate) Email Job Table
--------------------------------


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

--------------------------------
-- End: User (Candidate) Email Job Table
--------------------------------

--------------------------------
-- Start: Company Table
--------------------------------

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

--------------------------------
-- End: Company Table
--------------------------------

--------------------------------
-- Start: Recruiter Table
--------------------------------

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

--------------------------------
-- End: Recruiter Table
--------------------------------

--------------------------------
-- Start: Jobs Table
--------------------------------

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

--------------------------------
-- End: Jobs Table
--------------------------------

--------------------------------
-- Start: Views on User OAuth Tokens
--------------------------------

-- views on oauth tokens
create or replace view candidate_oauth_token with (security_invoker) as
select
  user_oauth_token.*
from user_oauth_token
inner join user_profile using (user_id);

create or replace view recruiter_oauth_token with (security_invoker) as
select
  user_oauth_token.*
from user_oauth_token
inner join recruiter using (user_id);

--------------------------------
-- End: Views on User OAuth Tokens
--------------------------------

--------------------------------
-- Start: Get Candidate/Recruiter Given Connected Email
--------------------------------

create or replace function get_user_profile_by_email (input text)
returns user_profile as 
$$
select
  user_profile.*
from public.user_profile
inner join public.user_oauth_token using (user_id)
where user_profile.email = input OR user_oauth_token.email = input
limit 1;
$$
language sql stable;

create or replace function get_recruiter_by_email (input text)
returns recruiter as 
$$
select
  recruiter.*
from public.recruiter
inner join public.user_oauth_token using (user_id)
where recruiter.email = input OR user_oauth_token.email = input
limit 1;
$$
language sql stable;

--------------------------------
-- End: Get Candidate/Recruiter Given Connected Email
--------------------------------

--------------------------------
-- Start: Recruiter Outbound Tables
--------------------------------


-- Recruiter outbound tables
-- These tables help us connect recruiter's existing outbound to SRC
create table public.recruiter_outbound_template (
  template_id uuid not null default uuid_generate_v4(),
  recruiter_id uuid references public.recruiter(user_id) on delete cascade not null,
  -- nullable job_id
  job_id uuid references public.job(job_id) on delete set null, 

  subject text not null,
  body text not null,
  normalized_content text not null,
  metadata jsonb not null default '{}'::jsonb,

  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone not null default now(),

  primary key (template_id)
);

create trigger handle_updated_at_recruiter_outbound_template before update on public.recruiter_outbound_template
  for each row execute procedure moddatetime (updated_at);

-- enable real-time
alter publication supabase_realtime add table public.recruiter_outbound_template;

-- enable RLS
alter table public.recruiter_outbound_template enable row level security;

create policy "Recruiters can view their outbound templates"
  on public.recruiter_outbound_template for select
  using ( auth.uid() = recruiter_id );

create policy "Recruiters can insert their outbound templates"
  on public.recruiter_outbound_template for insert
  with check ( auth.uid() = recruiter_id );

create policy "Recruiters can update their outbound templates"
  on public.recruiter_outbound_template for update
  using ( auth.uid() = recruiter_id );

create policy "Recruiters can delete their outbound templates"
  on public.recruiter_outbound_template for delete
  using ( auth.uid() = recruiter_id );

create or replace function list_similar_recruiter_outbound_templates (user_id uuid, input text)
returns table (
  template_id uuid,
  recruiter_id uuid,
  job_id uuid,
  subject text,
  body text,
  metadata jsonb,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  similarity real
) as 
$$
select
    recruiter_id,
    template_id,
    job_id,
    subject,
    body,
    metadata,
    created_at,
    updated_at,
    similarity(normalized_content, input) as "similarity"
from public.recruiter_outbound_template
where recruiter_id = user_id
and normalized_content % input
order by 9 desc;
$$
language sql stable;

-- add index for similarity search
create index idx_recruiter_outbound_template_normalized_content on public.recruiter_outbound_template using gin (normalized_content gin_trgm_ops);

-- default is 0.3
-- TODO: Adjust once we have more data
set pg_trgm.similarity_threshold = 0.5;

--------------------------------
--------------------------------

create table public.recruiter_outbound_message (
  recruiter_id uuid references public.recruiter(user_id) on delete cascade not null,
  -- message ID from the email provider
  message_id text not null,
  -- RFC2822 Message ID
  internal_message_id text not null,
  from_email text not null,
  to_email text not null,
  sent_at timestamp with time zone not null,

  template_id uuid references public.recruiter_outbound_template(template_id) on delete set null,

  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone not null default now(),

  primary key (message_id)
);

create trigger handle_updated_at_recruiter_outbound_message before update on public.recruiter_outbound_message
  for each row execute procedure moddatetime (updated_at);

-- enable real-time
alter publication supabase_realtime add table public.recruiter_outbound_message;

-- enable RLS
alter table public.recruiter_outbound_message enable row level security;

create policy "Recruiters can view their outbound messages"
  on public.recruiter_outbound_message for select
  using ( auth.uid() = recruiter_id );

--------------------------------
-- End: Recruiter Outbound Tables
--------------------------------

--------------------------------
-- Start: Candidate Company Inbound Tables & Triggers
--------------------------------

-- Candidate Company Inbound Tables
-- These tables help us generate a job board for candidates from recruiter outbound
create table public.candidate_company_inbound (
  candidate_email text not null,
  -- nullable
  candidate_id uuid references public.user_profile(user_id) on delete set null,
  company_id uuid references public.company(company_id) on delete cascade not null,
  recruiter_id uuid references public.recruiter(user_id) on delete set null,
  template_id uuid references public.recruiter_outbound_template(template_id) on delete cascade not null,
  -- nullable
  job_id uuid references public.job(job_id) on delete set null, 

  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone not null default now(),

  primary key (candidate_email, company_id, template_id)
);

create trigger handle_updated_at_candidate_company_inbound before update on public.candidate_company_inbound
  for each row execute procedure moddatetime (updated_at);

-- enable real-time
alter publication supabase_realtime add table public.candidate_company_inbound;

-- enable RLS
alter table public.candidate_company_inbound enable row level security;

-- TODO: Figure out company RLS policy
create policy "Recruiters can view their inbound candidates"
  on public.candidate_company_inbound for select
  using ( auth.uid() = recruiter_id );

create policy "Candidates can view their inbound"
  on public.candidate_company_inbound for select
  using ( auth.uid() = candidate_id );

create policy "Candidates can view inbound to their emails"
  on public.candidate_company_inbound for select
  using ( candidate_email in (
      select email
      from public.user_oauth_token
      where user_id = auth.uid()
  ));

-- create function to update job_id for a given template_id
create or replace function update_job_for_template_candidate_company_inbound_trigger()
returns trigger as
$$
begin
update public.candidate_company_inbound
set job_id = new.job_id
where template_id = new.template_id;
end;
$$
language plpgsql volatile;

-- create function to update candidate_id for a given candidate_email
create or replace function update_candidate_for_email_candidate_company_inbound_trigger()
returns trigger as
$$
begin
update public.candidate_company_inbound
set candidate_id = new.candidate_id
where candidate_email = new.candidate_email;
end;
$$
language plpgsql volatile;

-- create function to insert a row
create or replace function insert_candidate_company_inbound_trigger()
returns trigger as 
$$
begin
insert into public.candidate_company_inbound (
  candidate_email,
  company_id,
  template_id,
  recruiter_id
) 
select 
  new.candidate_email,
  new.template_id,
  new.recruiter_id,
  r.company_id
from new
inner join public.recruiter r on r.user_id = new.recruiter_id
on conflict do nothing;
end;
$$
language plpgsql volatile;

-- create a trigger for inserts into to recruiter_outbound_message
create trigger insert_candidate_company_inbound_trigger_after_insert after insert on public.recruiter_outbound_message
  for each row execute function insert_candidate_company_inbound_trigger();

-- create a trigger for inserts into to recruiter_outbound_template
create trigger update_job_for_template_candidate_company_inbound_trigger_insert_update after insert or update on public.recruiter_outbound_template
  for each row 
  when (new.job_id is not null)
  execute function update_job_for_template_candidate_company_inbound_trigger();

-- create a trigger for inserts into to user_oauth_token
create trigger update_candidate_for_email_candidate_company_inbound_trigger_user_oauth_token after insert or update on public.user_oauth_token
  for each row execute function update_candidate_for_email_candidate_company_inbound_trigger(); 

-- create a trigger for inserts into to user_profile
create trigger update_candidate_for_email_candidate_company_inbound_trigger_user_profile after insert on public.user_profile
  for each row execute function update_candidate_for_email_candidate_company_inbound_trigger();

--------------------------------
-- End: Candidate Company Inbound Tables & Triggers
--------------------------------
