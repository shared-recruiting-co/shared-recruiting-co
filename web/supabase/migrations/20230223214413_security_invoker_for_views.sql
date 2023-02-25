-- recreate views on oauth tokens with security invoker
create or replace view candidate_oauth_token with (security_invoker) as
select
  user_oauth_token.*
from user_oauth_token
inner join user_profile using (user_id);

create or replace view recruiter_oauth_token with (security_invoker) as
select
  user_oauth_token.*
from user_oauth_token
inner join recruiter using (user_id);
