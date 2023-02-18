alter table "public"."user_email_sync_history" drop constraint "user_email_sync_history_pkey";

drop index if exists "public"."user_email_sync_history_pkey";

alter table "public"."user_email_sync_history" alter column "email" set not null;

alter table "public"."user_email_sync_history" alter column "inbox_type" drop default;

create unique index user_email_sync_history_pkey on public.user_email_sync_history using btree (user_id, inbox_type, email);

alter table "public"."user_email_sync_history" add constraint "user_email_sync_history_pkey" primary key using index "user_email_sync_history_pkey";
