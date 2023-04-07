create type "public"."job_interest" as enum ('interested', 'not_interested', 'saved');

drop view if exists "public"."vw_job_board";

create table "public"."candidate_job_interest" (
    "candidate_id" uuid not null,
    "job_id" uuid not null,
    "interest" job_interest not null,
    "created_at" timestamp with time zone not null default now(),
    "updated_at" timestamp with time zone not null default now()
);


alter table "public"."candidate_job_interest" enable row level security;

CREATE UNIQUE INDEX candidate_job_interest_pkey ON public.candidate_job_interest USING btree (candidate_id, job_id);

alter table "public"."candidate_job_interest" add constraint "candidate_job_interest_pkey" PRIMARY KEY using index "candidate_job_interest_pkey";

alter table "public"."candidate_job_interest" add constraint "candidate_job_interest_candidate_id_fkey" FOREIGN KEY (candidate_id) REFERENCES user_profile(user_id) ON DELETE CASCADE not valid;

alter table "public"."candidate_job_interest" validate constraint "candidate_job_interest_candidate_id_fkey";

create or replace view vw_job_board with (security_invoker ) as (
  with verified_jobs as (
    select
      cc.candidate_id as user_id,
      cc.candidate_email as user_email,
      job.job_id,
      job.title as job_title,
      job.description_url as job_description_url,
      candidate_job_interest.interest as job_interest,
      company.company_name,
      company.website as company_website,
      recruiter.first_name || ' ' || recruiter.last_name as recruiter_name,
      recruiter.email as recruiter_email,
      cc.created_at as emailed_at,
      true as is_verified
    from public.candidate_company_inbound cc
    inner join public.company on cc.company_id = company.company_id
    inner join public.recruiter on cc.recruiter_id = recruiter.user_id
    inner join public.job on cc.job_id = job.job_id
    left join public.candidate_job_interest on cc.candidate_id = candidate_job_interest.candidate_id and cc.job_id = candidate_job_interest.job_id
  ),
  unofficial_jobs as (
    select
      user_id,
      user_email,
      user_email_job.job_id,
      job_title,
      '' as job_description_url,
      candidate_job_interest.interest as job_interest,
      company as company_name,
      '' as company_website,
      data ->> 'recruiter' as recruiter_name,
      data ->> 'recruiter_email' as recruiter_email,
      emailed_at,
      false as is_verified
    from public.user_email_job
    left join public.candidate_job_interest
    on user_email_job.user_id = candidate_job_interest.candidate_id
    and user_email_job.job_id = candidate_job_interest.job_id
  )
  -- line up columns
  select * from verified_jobs
  union all
  select * from unofficial_jobs
);


create policy "Candidates can insert their own job interest"
on "public"."candidate_job_interest"
as permissive
for insert
to public
with check ((auth.uid() = candidate_id));


create policy "Candidates can update their own job interest"
on "public"."candidate_job_interest"
as permissive
for update
to public
using ((auth.uid() = candidate_id));


create policy "Candidates can view their own job interest"
on "public"."candidate_job_interest"
as permissive
for select
to public
using ((auth.uid() = candidate_id));


CREATE TRIGGER handle_updated_at_candidate_job_interest BEFORE UPDATE ON public.candidate_job_interest FOR EACH ROW EXECUTE FUNCTION moddatetime('updated_at');


