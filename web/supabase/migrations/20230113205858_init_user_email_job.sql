-- This script was generated by the Schema Diff utility in pgAdmin 4
-- For the circular dependencies, the order in which Schema Diff writes the objects is not very sophisticated
-- and may require manual changes to the script to ensure changes are applied in the correct order.
-- Please report an issue for any failure with the reproduction steps.

CREATE TABLE IF NOT EXISTS public.user_email_job
(
    job_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    user_email text COLLATE pg_catalog."default" NOT NULL,
    email_thread_id text COLLATE pg_catalog."default" NOT NULL,
    emailed_at timestamp with time zone NOT NULL,
    company text COLLATE pg_catalog."default" NOT NULL,
    job_title text COLLATE pg_catalog."default" NOT NULL,
    data jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT user_email_job_pkey PRIMARY KEY (job_id),
    CONSTRAINT user_email_job_user_email_email_thread_id_key UNIQUE (user_email, email_thread_id),
    CONSTRAINT user_email_job_user_id_company_job_title_key UNIQUE (user_id, company, job_title),
    CONSTRAINT user_email_job_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES auth.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.user_email_job
    OWNER to postgres;

ALTER TABLE IF EXISTS public.user_email_job
    ENABLE ROW LEVEL SECURITY;

GRANT ALL ON TABLE public.user_email_job TO anon;

GRANT ALL ON TABLE public.user_email_job TO authenticated;

GRANT ALL ON TABLE public.user_email_job TO postgres;

GRANT ALL ON TABLE public.user_email_job TO service_role;
CREATE POLICY "Users can view their own jobs"
    ON public.user_email_job
    AS PERMISSIVE
    FOR SELECT
    TO public
    USING ((auth.uid() = user_id));

CREATE TRIGGER handle_updated_at_user_email_job
    BEFORE UPDATE 
    ON public.user_email_job
    FOR EACH ROW
    EXECUTE FUNCTION extensions.moddatetime('updated_at');

-- manually added
ALTER PUBLICATION supabase_realtime ADD TABLE user_email_job;
