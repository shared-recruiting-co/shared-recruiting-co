create or replace view candidate_job_count_unverified with (security_invoker) as
select
  user_id,
  count(distinct job_id) as num_jobs
from user_email_job
group by user_id;
