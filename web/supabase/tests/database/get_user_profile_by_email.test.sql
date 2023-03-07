--------------
-- Get User Profile by Email Test
--------------
-- Test the get_user_profile_by_email helper function
--------------
begin;
select plan( 3 );

\set user_id 'f50bc4c9-b14e-406b-8b47-e1d0503f4ff9'
\set user_email 'candidate@example.com'
\set user_email_two 'personal@example.com'
\set first_name 'Test'
\set last_name 'User'

-- insert user
insert into auth.users (id, email, raw_app_meta_data, raw_user_meta_data, aud, role)
values (:'user_id', :'user_email', '{"provider":"email"}', '{"full_name":"Test User"}', 'authenticated', 'authenticated');
-- insert user_profile
insert into user_profile (user_id, email, first_name, last_name) values (:'user_id', :'user_email', :'first_name', :'last_name');
-- CHECK: get user profile by email (no oauth token)
select is(
    (select user_id from get_user_profile_by_email(:'user_email')),
    :'user_id',
  'user profile with no oauth token fetched correctly'
);
-- insert user_oauth_token with same user_profile email
insert into user_oauth_token (user_id, email, provider, token) values (:'user_id', :'user_email', 'google', '{}'::jsonb);
-- CHECK: get user profile by email
select is(
    (select user_id from get_user_profile_by_email(:'user_email')),
    :'user_id',
  'user profile with same email oauth token fetched correctly'
);

-- insert user oauth token with different email
insert into user_oauth_token (user_id, email, provider, token) values (:'user_id', :'user_email_two', 'google', '{}'::jsonb);
-- CHECK: get user profile by email
select is(
    (select user_id from get_user_profile_by_email(:'user_email_two')),
    :'user_id',
  'user profile with different email oauth token fetched correctly'
);


select * from finish();
rollback;
