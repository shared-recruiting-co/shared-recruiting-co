create unique index unique_email_provider on public.user_oauth_token using btree (email, provider);

alter table "public"."user_oauth_token" add constraint "unique_email_provider" unique using index "unique_email_provider";
