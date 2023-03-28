create or replace view public.outbound_template_recipient_count with (security_invoker) as
select 
   template_id,
   count(distinct to_email) as num_recipients
from recruiter_outbound_message
group by template_id;

create or replace view public.job_candidate_count with (security_invoker) as
select
  job_id,
  count(distinct candidate_email) as num_candidates
from candidate_company_inbound
group by job_id;
