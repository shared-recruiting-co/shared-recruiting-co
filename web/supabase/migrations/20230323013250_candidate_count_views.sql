create or replace view "public"."job_candidate_count" as  SELECT candidate_company_inbound.job_id,
    count(DISTINCT candidate_company_inbound.candidate_email) AS num_candidates
   FROM candidate_company_inbound
  GROUP BY candidate_company_inbound.job_id;


create or replace view "public"."outbound_template_recipient_count" as  SELECT recruiter_outbound_message.template_id,
    count(DISTINCT recruiter_outbound_message.to_email) AS num_recipients
   FROM recruiter_outbound_message
  GROUP BY recruiter_outbound_message.template_id;



