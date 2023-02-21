CREATE UNIQUE INDEX unique_email_provider ON public.user_oauth_token USING btree (email, provider);

alter table "public"."user_oauth_token" add constraint "unique_email_provider" UNIQUE using index "unique_email_provider";


