create or replace view "public"."candidate_job_count" with (security_invoker) as  
  select candidate_company_inbound.candidate_id,
  count(distinct candidate_company_inbound.job_id) as num_jobs
from candidate_company_inbound
where candidate_company_inbound.candidate_id is not null
group by candidate_company_inbound.candidate_id;
