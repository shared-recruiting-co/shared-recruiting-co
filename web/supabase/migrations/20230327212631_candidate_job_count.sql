create or replace view "public"."candidate_job_count" as  SELECT candidate_company_inbound.candidate_id,
    count(DISTINCT candidate_company_inbound.job_id) AS num_jobs
   FROM candidate_company_inbound
  GROUP BY candidate_company_inbound.candidate_id;
