--------------
-- Candidate Company Inbound
--------------
-- This should always pass. If it doesn't the test infra is broken.
--------------
begin;
select plan( 6 );

-- declare variables
\set company_id 'f50bc4c9-b14e-406b-8b47-e1d0503f4ff9'
\set recruiter_id 'ad5ac50b-b8e6-4d7c-848f-678b7eb07fcb'
\set template_id '695566d2-1af9-4ea7-af85-b488b6db5300'
\set job_id 'ac1d6268-94a1-4b7a-a11e-4290b9b4cedf'
\set message_id 'A'
\set recruiter_email 'recruiter@example.com'
\set candidate_id 'ef6a6548-dec6-4db6-a929-6e2fd96b7464'
\set candidate_email 'candidate@example.com'

insert into auth.users (id, email, raw_app_meta_data, raw_user_meta_data, aud, role)
values (:'recruiter_id', :'recruiter_email', '{"provider":"email"}', '{"full_name":"Test User"}', 'authenticated', 'authenticated');

-- 1. create a company
insert into company (company_id, company_name, website) values (:'company_id', 'Test Company', 'https://test.com');
-- 2. create a recruiter
insert into recruiter (user_id, email, first_name, last_name, company_id) values (:'recruiter_id', :'recruiter_email', 'Test', 'User', :'company_id');

-- 3. create a outbound template
insert into recruiter_outbound_template (template_id, recruiter_id, subject, body, normalized_content) values (:'template_id', :'recruiter_id', 'Test Subject', 'Test Body', 'Test Subject Test Body');

-- 4. create a outbound message
insert into recruiter_outbound_message (message_id, template_id, recruiter_id, internal_message_id, from_email, to_email, sent_at) values (:'message_id', :'template_id', :'recruiter_id',  '123', :'recruiter_email', :'candidate_email', now());

-- 5. CHECK candidate_company_inbound
select is(
  (select count(*) from candidate_company_inbound),
  1::bigint,
  'candidate company inbound row created'
);

SELECT is(
    (
    select (
        candidate_email, 
        candidate_id, 
        company_id, 
        recruiter_id, 
        template_id, 
        job_id,
        now(),
        now()
    )::candidate_company_inbound
    from candidate_company_inbound
    ),
    row( 
        :'candidate_email', 
        null, 
        :'company_id', 
        :'recruiter_id', 
        :'template_id', 
        null,
        now(),
        now()
    )::candidate_company_inbound,
    'candidate company inbound row created correctly'
);

-- 6. create a job
insert into job (job_id, company_id, recruiter_id, title, description_url) values (:'job_id', :'company_id', :'recruiter_id', 'Test Job', 'https://test.com');
-- 7. update template's job
update recruiter_outbound_template set job_id = :'job_id' where template_id = :'template_id';
-- 8. CHECK candidate_company_inbound
SELECT is(
    (
    select (
        candidate_email, 
        candidate_id, 
        company_id, 
        recruiter_id, 
        template_id, 
        job_id,
        now(),
        now()
    )::candidate_company_inbound
    from candidate_company_inbound
    ),
    row( 
        :'candidate_email', 
        null, 
        :'company_id', 
        :'recruiter_id', 
        :'template_id', 
        :'job_id',
        now(),
        now()
    )::candidate_company_inbound,
    'candidate company inbound row updated with job'
);
-- 9. create a user_profile
insert into auth.users (id, email, raw_app_meta_data, raw_user_meta_data, aud, role)
values (:'candidate_id', :'candidate_email', '{"provider":"email"}', '{"full_name":"Test User"}', 'authenticated', 'authenticated');
insert into user_profile (user_id, email, first_name, last_name) values (:'candidate_id', :'candidate_email', 'Test', 'User');

-- 10. CHECK candidate_company_inbound
SELECT is(
    (
    select (
        candidate_email, 
        candidate_id, 
        company_id, 
        recruiter_id, 
        template_id, 
        job_id,
        now(),
        now()
    )::candidate_company_inbound
    from candidate_company_inbound
    ),
    row( 
        :'candidate_email', 
        :'candidate_id',
        :'company_id', 
        :'recruiter_id', 
        :'template_id', 
        :'job_id',
        now(),
        now()
    )::candidate_company_inbound,
    'candidate company inbound row updated with candidate'
);

-- verify candidate_id is inserted if it exists at time of outbound message

-- 11. add another outbound template
\set template_id_2 'f75b017d-484f-4704-b6a3-1e185fdd2ee4'
insert into recruiter_outbound_template (template_id, recruiter_id, subject, body, normalized_content) values (:'template_id_2', :'recruiter_id', 'Test Subject', 'Test Body', 'Test Subject Test Body');
-- 12. add another outbound message
\set message_id_2 'B'
insert into recruiter_outbound_message (message_id, template_id, recruiter_id, internal_message_id, from_email, to_email, sent_at) values (:'message_id_2', :'template_id_2', :'recruiter_id',  '123', :'recruiter_email', :'candidate_email', now());

-- 13. CHECK candidate_company_inbound
SELECT is(
    (
    select (
        candidate_email, 
        candidate_id, 
        company_id, 
        recruiter_id, 
        template_id, 
        job_id,
        now(),
        now()
    )::candidate_company_inbound
    from candidate_company_inbound
    where template_id = :'template_id_2'
    ),
    row( 
        :'candidate_email', 
        :'candidate_id',
        :'company_id', 
        :'recruiter_id', 
        :'template_id_2', 
        null,
        now(),
        now()
    )::candidate_company_inbound,
    'candidate company inbound row updated with candidate'
);

-- verify insert/updates to user_oauth_token
\set candidate_email_2 'two@candidate.com'

-- 12. add another outbound message
\set message_id_3 'C'
insert into recruiter_outbound_message (message_id, template_id, recruiter_id, internal_message_id, from_email, to_email, sent_at) values (:'message_id_3', :'template_id', :'recruiter_id',  '123', :'recruiter_email', :'candidate_email_2', now());

-- insert a user_oauth_token
insert into user_oauth_token (user_id, email, provider, token) values (:'candidate_id', :'candidate_email_2', 'google', '{}'::jsonb);

-- 13. CHECK candidate_company_inbound
SELECT is(
    (
    select (
        candidate_email, 
        candidate_id, 
        company_id, 
        recruiter_id, 
        template_id, 
        job_id,
        now(),
        now()
    )::candidate_company_inbound
    from candidate_company_inbound
    where template_id = :'template_id'
    and candidate_email = :'candidate_email_2'
    ),
    row( 
        :'candidate_email_2', 
        :'candidate_id',
        :'company_id', 
        :'recruiter_id', 
        :'template_id', 
        :'job_id',
        now(),
        now()
    )::candidate_company_inbound,
    'candidate company inbound row updated with candidate'
);

select * from finish();

rollback;
