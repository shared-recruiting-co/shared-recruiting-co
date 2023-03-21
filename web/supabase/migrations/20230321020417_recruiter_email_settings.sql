alter table "public"."recruiter" add column "email_settings" jsonb not null default '{}'::jsonb;


