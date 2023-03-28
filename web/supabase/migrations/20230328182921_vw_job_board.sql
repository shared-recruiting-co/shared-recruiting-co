create or replace view public.vw_job_board with (security_invoker ) as (
  with verified_jobs as (
    select 
      cc.candidate_id as user_id,
      cc.candidate_email as user_email,
      job.job_id,
      job.title as job_title,
      job.description_url as job_description_url,
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
  ),
  unofficial_jobs as (
    select 
      user_id,
      user_email,
      job_id,
      job_title,
      '' as job_description_url,
      company as company_name,
      '' as company_website,
      data ->> 'recruiter' as recruiter_name,
      data ->> 'recruiter_email' as recruiter_email,
      emailed_at,
      false as is_verified
    from public.user_email_job
  )
  -- line up columns
  select * from verified_jobs
  union all
  select * from unofficial_jobs
);
