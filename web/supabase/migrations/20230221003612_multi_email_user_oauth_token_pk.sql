alter table "public"."user_oauth_token" drop constraint "user_oauth_token_pkey";

drop index if exists "public"."user_oauth_token_pkey";

alter table "public"."user_oauth_token" alter column "email" set not null;

CREATE UNIQUE INDEX user_oauth_token_pkey ON public.user_oauth_token USING btree (user_id, email, provider);

alter table "public"."user_oauth_token" add constraint "user_oauth_token_pkey" PRIMARY KEY using index "user_oauth_token_pkey";


