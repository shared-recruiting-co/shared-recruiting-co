create extension if not exists "pgtap" with schema "extensions";


create table "public"."candidate_company_inbound" (
    "candidate_email" text not null,
    "candidate_id" uuid,
    "company_id" uuid not null,
    "recruiter_id" uuid,
    "template_id" uuid not null,
    "job_id" uuid,
    "created_at" timestamp with time zone not null default now(),
    "updated_at" timestamp with time zone not null default now()
);


alter table "public"."candidate_company_inbound" enable row level security;

CREATE UNIQUE INDEX candidate_company_inbound_pkey ON public.candidate_company_inbound USING btree (candidate_email, company_id, template_id);

alter table "public"."candidate_company_inbound" add constraint "candidate_company_inbound_pkey" PRIMARY KEY using index "candidate_company_inbound_pkey";

alter table "public"."candidate_company_inbound" add constraint "candidate_company_inbound_candidate_id_fkey" FOREIGN KEY (candidate_id) REFERENCES user_profile(user_id) ON DELETE SET NULL not valid;

alter table "public"."candidate_company_inbound" validate constraint "candidate_company_inbound_candidate_id_fkey";

alter table "public"."candidate_company_inbound" add constraint "candidate_company_inbound_company_id_fkey" FOREIGN KEY (company_id) REFERENCES company(company_id) ON DELETE CASCADE not valid;

alter table "public"."candidate_company_inbound" validate constraint "candidate_company_inbound_company_id_fkey";

alter table "public"."candidate_company_inbound" add constraint "candidate_company_inbound_job_id_fkey" FOREIGN KEY (job_id) REFERENCES job(job_id) ON DELETE SET NULL not valid;

alter table "public"."candidate_company_inbound" validate constraint "candidate_company_inbound_job_id_fkey";

alter table "public"."candidate_company_inbound" add constraint "candidate_company_inbound_recruiter_id_fkey" FOREIGN KEY (recruiter_id) REFERENCES recruiter(user_id) ON DELETE SET NULL not valid;

alter table "public"."candidate_company_inbound" validate constraint "candidate_company_inbound_recruiter_id_fkey";

alter table "public"."candidate_company_inbound" add constraint "candidate_company_inbound_template_id_fkey" FOREIGN KEY (template_id) REFERENCES recruiter_outbound_template(template_id) ON DELETE CASCADE not valid;

alter table "public"."candidate_company_inbound" validate constraint "candidate_company_inbound_template_id_fkey";

set check_function_bodies = off;

CREATE OR REPLACE FUNCTION public.insert_candidate_company_inbound_trigger()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
begin
with input as (
  select 
    *
  from (values (new.to_email, new.recruiter_id, new.template_id)) as t (candidate_email, recruiter_id, template_id)
), 
c as (
  select
    user_id as candidate_id,
    coalesce(email, new.to_email) as candidate_email
  from get_user_profile_by_email(new.to_email)
),
t as (
  select 
    template_id,
    job_id
  from public.recruiter_outbound_template
  where template_id = new.template_id
)
insert into public.candidate_company_inbound (
  candidate_email,
  candidate_id,
  template_id,
  job_id,
  recruiter_id,
  company_id
) 
select 
  input.candidate_email,
  c.candidate_id,
  input.template_id,
  t.job_id,
  input.recruiter_id,
  r.company_id
from input
inner join public.recruiter r on r.user_id = input.recruiter_id
inner join t using (template_id)
left join c using (candidate_email)
on conflict do nothing;
return null;
end;
$function$
;

CREATE OR REPLACE FUNCTION public.update_candidate_for_email_candidate_company_inbound_trigger()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
begin
update public.candidate_company_inbound
set candidate_id = new.user_id
where candidate_email = new.email;
return null;
end;
$function$
;

CREATE OR REPLACE FUNCTION public.update_job_for_template_candidate_company_inbound_trigger()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
begin
update public.candidate_company_inbound
set job_id = new.job_id
where template_id = new.template_id;
return null;
end;
$function$
;

CREATE OR REPLACE FUNCTION public.get_recruiter_by_email(input text)
 RETURNS recruiter
 LANGUAGE sql
 STABLE
AS $function$
select
  recruiter.*
from public.recruiter
left join public.user_oauth_token using (user_id)
where recruiter.email = input OR user_oauth_token.email = input
limit 1;
$function$
;

CREATE OR REPLACE FUNCTION public.get_user_profile_by_email(input text)
 RETURNS user_profile
 LANGUAGE sql
 STABLE
AS $function$
select
  user_profile.*
from public.user_profile
left join public.user_oauth_token using (user_id)
where user_profile.email = input OR user_oauth_token.email = input
limit 1;
$function$
;

create policy "Candidates can view inbound to their emails"
on "public"."candidate_company_inbound"
as permissive
for select
to public
using ((candidate_email IN ( SELECT user_oauth_token.email
   FROM user_oauth_token
  WHERE (user_oauth_token.user_id = auth.uid()))));


create policy "Candidates can view their inbound"
on "public"."candidate_company_inbound"
as permissive
for select
to public
using ((auth.uid() = candidate_id));


create policy "Recruiters can view their inbound candidates"
on "public"."candidate_company_inbound"
as permissive
for select
to public
using ((auth.uid() = recruiter_id));


CREATE TRIGGER handle_updated_at_candidate_company_inbound BEFORE UPDATE ON public.candidate_company_inbound FOR EACH ROW EXECUTE FUNCTION moddatetime('updated_at');

CREATE TRIGGER insert_candidate_company_inbound_trigger_after_insert AFTER INSERT ON public.recruiter_outbound_message FOR EACH ROW EXECUTE FUNCTION insert_candidate_company_inbound_trigger();

CREATE TRIGGER candidate_company_inbound_trigger_recruiter_outbound_template AFTER INSERT OR UPDATE ON public.recruiter_outbound_template FOR EACH ROW WHEN ((new.job_id IS NOT NULL)) EXECUTE FUNCTION update_job_for_template_candidate_company_inbound_trigger();

CREATE TRIGGER candidate_company_inbound_trigger_user_oauth_token AFTER INSERT OR UPDATE ON public.user_oauth_token FOR EACH ROW EXECUTE FUNCTION update_candidate_for_email_candidate_company_inbound_trigger();

CREATE TRIGGER candidate_company_inbound_trigger_user_profile AFTER INSERT ON public.user_profile FOR EACH ROW EXECUTE FUNCTION update_candidate_for_email_candidate_company_inbound_trigger();
