create extension if not exists "pg_trgm" with schema "extensions";


create table "public"."recruiter_outbound_message" (
    "recruiter_id" uuid not null,
    "message_id" text not null,
    "internal_message_id" text not null,
    "from_email" text not null,
    "to_email" text not null,
    "sent_at" timestamp with time zone not null,
    "template_id" uuid,
    "created_at" timestamp with time zone not null default now(),
    "updated_at" timestamp with time zone not null default now()
);


alter table "public"."recruiter_outbound_message" enable row level security;

create table "public"."recruiter_outbound_template" (
    "template_id" uuid not null default uuid_generate_v4(),
    "recruiter_id" uuid not null,
    "job_id" uuid,
    "subject" text not null,
    "body" text not null,
    "normalized_content" text not null,
    "metadata" jsonb not null default '{}'::jsonb,
    "created_at" timestamp with time zone not null default now(),
    "updated_at" timestamp with time zone not null default now()
);


alter table "public"."recruiter_outbound_template" enable row level security;

CREATE INDEX idx_recruiter_outbound_template_normalized_content ON public.recruiter_outbound_template USING gin (normalized_content gin_trgm_ops);

CREATE UNIQUE INDEX recruiter_outbound_message_pkey ON public.recruiter_outbound_message USING btree (message_id);

CREATE UNIQUE INDEX recruiter_outbound_template_pkey ON public.recruiter_outbound_template USING btree (template_id);

alter table "public"."recruiter_outbound_message" add constraint "recruiter_outbound_message_pkey" PRIMARY KEY using index "recruiter_outbound_message_pkey";

alter table "public"."recruiter_outbound_template" add constraint "recruiter_outbound_template_pkey" PRIMARY KEY using index "recruiter_outbound_template_pkey";

alter table "public"."recruiter_outbound_message" add constraint "recruiter_outbound_message_recruiter_id_fkey" FOREIGN KEY (recruiter_id) REFERENCES recruiter(user_id) ON DELETE CASCADE not valid;

alter table "public"."recruiter_outbound_message" validate constraint "recruiter_outbound_message_recruiter_id_fkey";

alter table "public"."recruiter_outbound_message" add constraint "recruiter_outbound_message_template_id_fkey" FOREIGN KEY (template_id) REFERENCES recruiter_outbound_template(template_id) ON DELETE SET NULL not valid;

alter table "public"."recruiter_outbound_message" validate constraint "recruiter_outbound_message_template_id_fkey";

alter table "public"."recruiter_outbound_template" add constraint "recruiter_outbound_template_job_id_fkey" FOREIGN KEY (job_id) REFERENCES job(job_id) ON DELETE SET NULL not valid;

alter table "public"."recruiter_outbound_template" validate constraint "recruiter_outbound_template_job_id_fkey";

alter table "public"."recruiter_outbound_template" add constraint "recruiter_outbound_template_recruiter_id_fkey" FOREIGN KEY (recruiter_id) REFERENCES recruiter(user_id) ON DELETE CASCADE not valid;

alter table "public"."recruiter_outbound_template" validate constraint "recruiter_outbound_template_recruiter_id_fkey";

set check_function_bodies = off;

CREATE OR REPLACE FUNCTION public.list_similar_recruiter_outbound_templates(user_id uuid, input text)
 RETURNS TABLE(template_id uuid, recruiter_id uuid, job_id uuid, subject text, body text, metadata jsonb, created_at timestamp with time zone, updated_at timestamp with time zone, similarity real)
 LANGUAGE sql
 STABLE
AS $function$
select
    template_id,
    recruiter_id,
    job_id,
    subject,
    body,
    metadata,
    created_at,
    updated_at,
    similarity(normalized_content, input) as "similarity"
from public.recruiter_outbound_template
where recruiter_id = user_id
and normalized_content % input
order by 9 desc;
$function$
;

create policy "Recruiters can view their outbound messages"
on "public"."recruiter_outbound_message"
as permissive
for select
to public
using ((auth.uid() = recruiter_id));


create policy "Recruiters can delete their outbound templates"
on "public"."recruiter_outbound_template"
as permissive
for delete
to public
using ((auth.uid() = recruiter_id));


create policy "Recruiters can insert their outbound templates"
on "public"."recruiter_outbound_template"
as permissive
for insert
to public
with check ((auth.uid() = recruiter_id));


create policy "Recruiters can update their outbound templates"
on "public"."recruiter_outbound_template"
as permissive
for update
to public
using ((auth.uid() = recruiter_id));


create policy "Recruiters can view their outbound templates"
on "public"."recruiter_outbound_template"
as permissive
for select
to public
using ((auth.uid() = recruiter_id));


CREATE TRIGGER handle_updated_at_recruiter_outbound_message BEFORE UPDATE ON public.recruiter_outbound_message FOR EACH ROW EXECUTE FUNCTION moddatetime('updated_at');

CREATE TRIGGER handle_updated_at_recruiter_outbound_template BEFORE UPDATE ON public.recruiter_outbound_template FOR EACH ROW EXECUTE FUNCTION moddatetime('updated_at');

set pg_trgm.similarity_threshold = 0.5;
