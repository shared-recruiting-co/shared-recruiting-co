-- This script was generated by the Schema Diff utility in pgAdmin 4
-- For the circular dependencies, the order in which Schema Diff writes the objects is not very sophisticated
-- and may require manual changes to the script to ensure changes are applied in the correct order.
-- Please report an issue for any failure with the reproduction steps.

CREATE OR REPLACE VIEW public.recruiter_oauth_token
 AS
 SELECT 
    user_oauth_token.*
   FROM user_oauth_token
     JOIN recruiter USING (user_id);

ALTER TABLE public.recruiter_oauth_token
    OWNER TO postgres;

GRANT ALL ON TABLE public.recruiter_oauth_token TO authenticated;
GRANT ALL ON TABLE public.recruiter_oauth_token TO postgres;
GRANT ALL ON TABLE public.recruiter_oauth_token TO anon;
GRANT ALL ON TABLE public.recruiter_oauth_token TO service_role;

CREATE OR REPLACE VIEW public.candidate_oauth_token
 AS
 SELECT
    user_oauth_token.*
   FROM user_oauth_token
     JOIN user_profile USING (user_id);

ALTER TABLE public.candidate_oauth_token
    OWNER TO postgres;

GRANT ALL ON TABLE public.candidate_oauth_token TO authenticated;
GRANT ALL ON TABLE public.candidate_oauth_token TO postgres;
GRANT ALL ON TABLE public.candidate_oauth_token TO anon;
GRANT ALL ON TABLE public.candidate_oauth_token TO service_role;
