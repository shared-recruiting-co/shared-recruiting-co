--------------
-- Canary Test
--------------
-- This should always pass. If it doesn't the test infra is broken.
--------------
begin;
select plan( 1 );

select results_eq(
    'select * from (values (1))  as t',
    $$VALUES (1)$$,
    'canary'
);


select * from finish();
rollback;
