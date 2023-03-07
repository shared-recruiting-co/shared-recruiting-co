--------------
-- Canary Test
--------------
-- This should always pass. If it doesn't the test infra is broken.
--------------
begin;
select plan( 2 );

-- is pgtap installed
select results_eq(
    'select * from (values (1))  as t',
    $$VALUES (1)$$,
    'canary'
);

-- is supabase database setup
SELECT has_column(
    'auth',
    'users',
    'id',
    'id should exist'
);

select * from finish();
rollback;
