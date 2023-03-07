--------------
-- Get Recruiter by Email Test
--------------
-- Test the get_recruiter_by_email helper function
--------------
begin;
select plan( 3 );

\set company_id 'f50bc4c9-b14e-406b-8b47-e1d0503f4ff9'
\set user_id 'f50bc4c9-b14e-406b-8b47-e1d0503f4ff9'
\set user_email 'recruiter@example.com'
\set user_email_two 'personal@example.com'
\set first_name 'Test'
\set last_name 'User'

-- insert user
insert into auth.users (id, email, raw_app_meta_data, raw_user_meta_data, aud, role)
values (:'user_id', :'user_email', '{"provider":"email"}', '{"full_name":"Test User"}', 'authenticated', 'authenticated');
-- insert company and recruiter
insert into company (company_id, company_name, website) values (:'company_id', 'Test Company', 'https://test.com');
insert into recruiter (user_id, email, first_name, last_name, company_id) values (:'user_id', :'user_email', :'first_name', :'last_name', :'company_id');
-- CHECK: get recruiter by email (no oauth token)
select is(
    (select user_id from get_recruiter_by_email(:'user_email')),
    :'user_id',
  'recruiter with no oauth token fetched correctly'
);
-- insert user_oauth_token with same recruiter email
insert into user_oauth_token (user_id, email, provider, token) values (:'user_id', :'user_email', 'google', '{}'::jsonb);
-- CHECK: get recruiter by email
select is(
    (select user_id from get_recruiter_by_email(:'user_email')),
    :'user_id',
  'recruiter with same email oauth token fetched correctly'
);

-- insert user oauth token with different email
insert into user_oauth_token (user_id, email, provider, token) values (:'user_id', :'user_email_two', 'google', '{}'::jsonb);
-- CHECK: get recruiter by email
select is(
    (select user_id from get_recruiter_by_email(:'user_email_two')),
    :'user_id',
  'recruiter with different email oauth token fetched correctly'
);


select * from finish();
rollback;
