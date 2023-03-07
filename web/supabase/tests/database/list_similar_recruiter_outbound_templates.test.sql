--------------
-- List Similar Recruiter Outbound Templates Test
--------------
-- Test the list_similar_recruiter_outbound_templates helper function
-- The function fetches the top n similar recruiter outbound templates to a given input
--------------
begin;
select plan( 3 );

\set company_id 'f50bc4c9-b14e-406b-8b47-e1d0503f4ff9'
\set user_id 'ef6a6548-dec6-4db6-a929-6e2fd96b7464'
\set user_email 'recruiter@example.com'
\set template_id '695566d2-1af9-4ea7-af85-b488b6db5300'
\set job_id 'ac1d6268-94a1-4b7a-a11e-4290b9b4cedf'
\set first_name 'Test'
\set last_name 'User'
\set template_subject E'Hey, How are things with you?'
\set template_body E'Hey {{first_name}}, I hope this finds you doing well. (= I had reached out to you earlier regarding a few opportunities. I thought I\'d check in with you and see how things are going. I am guessing you still may not be actively looking, but there may be some exciting possibilities worth exploring. Chances are we could have a match for you based on your impressive background. There are many companies that would like to connect to an engineer with your expertise and I\'d like to share whatever opportunities match your interests most. Let me know how things are going and the best time and number to reach you.'
\set template_normalized_content :'template_subject' || ' ' || :'template_body'

-- 1. insert user
insert into auth.users (id, email, raw_app_meta_data, raw_user_meta_data, aud, role)
values (:'user_id', :'user_email', '{"provider":"email"}', '{"full_name":"Test User"}', 'authenticated', 'authenticated');
-- 2. insert company and recruiter
insert into company (company_id, company_name, website) values (:'company_id', 'Test Company', 'https://test.com');
insert into recruiter (user_id, email, first_name, last_name, company_id) values (:'user_id', :'user_email', :'first_name', :'last_name', :'company_id');
-- 3. insert a job
insert into job (job_id, company_id, recruiter_id, title, description_url) values (:'job_id', :'company_id', :'user_id', 'Test Job', 'https://test.com');
-- 4. insert an outbound template
insert into recruiter_outbound_template (template_id, recruiter_id, subject, body, normalized_content) values (:'template_id', :'user_id', :'template_subject', :'template_body', :'template_normalized_content');

-- CHECK: list_similar_recruiter_outbound_templates
select is(
    (select
        template_id 
        from list_similar_recruiter_outbound_templates(:'user_id', :'template_normalized_content')
    ),
    :'template_id',
  'identical content should return the same template'
);

\set different_content E'Gartner forecast "by 2026, at least 60% of organizations procuring mission-critical software solutions will mandate SBOMs disclosures," yet many product leaders do not have an accurate view of their software supply chain. The report addresses several critical insights that will empower product leaders to understand and better address the challenges and opportunities presented across this high-priority emerging trend. In this report you will learn:'

select is(
    (select
        template_id 
        from list_similar_recruiter_outbound_templates(:'user_id', :'different_content')
    ),
    null,
  'different content should not return a template'
);

\set similar_content :'template_subject' || ' ' || E'Hey John, I hope this finds you doing well. (= I had reached out to you earlier regarding a few opportunities. I thought I\'d check in with you and see how things are going. I am guessing you still may not be actively looking, but there may be some exciting possibilities worth exploring. Chances are we could have a match for you based on your impressive background. There are many companies that would like to connect to an engineer with your expertise and I\'d like to share whatever opportunities match your interests most. Let me know how things are going and the best time and number to reach you.'

select is(
    (select
        template_id 
        from list_similar_recruiter_outbound_templates(:'user_id', :'similar_content')
    ),
    :'template_id',
  'similar content should return a template'
);

select * from finish();
rollback;
