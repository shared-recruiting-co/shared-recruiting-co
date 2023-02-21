alter table "public"."user_oauth_token" drop constraint "user_oauth_token_pkey";

drop index if exists "public"."user_oauth_token_pkey";

alter table "public"."user_oauth_token" alter column "email" set not null;

CREATE UNIQUE INDEX user_oauth_token_pkey ON public.user_oauth_token USING btree (user_id, email, provider);

alter table "public"."user_oauth_token" add constraint "user_oauth_token_pkey" PRIMARY KEY using index "user_oauth_token_pkey";

create or replace function get_user_profile_by_email (input text)
returns user_profile as 
$$
select
  user_profile.*
from public.user_profile
inner join public.user_oauth_token using (user_id)
where user_profile.email = input OR user_oauth_token.email = input;
$$
language sql stable;

create or replace function get_recruiter_by_email (input text)
returns recruiter as 
$$
select
  recruiter.*
from public.recruiter
inner join public.user_oauth_token using (user_id)
where recruiter.email = input OR user_oauth_token.email = input;
$$
language sql stable;
