create or replace procedure auth.login_as_user(user_email text)
    language plpgsql
    as $$
declare
    auth_user auth.users;
begin
    select
        * into auth_user
    from
        auth.users
    where
        email = user_email;
    execute format('set request.jwt.claim.sub=%L', (auth_user).id::text);
    execute format('set request.jwt.claim.role=%I', (auth_user).role);
    execute format('set request.jwt.claim.email=%L', (auth_user).email);
    execute format('set request.jwt.claims=%L', json_strip_nulls(json_build_object('app_metadata', (auth_user).raw_app_meta_data))::text);
    raise notice '%', format( 'set role %I; -- logging in as %L (%L)', (auth_user).role, (auth_user).id, (auth_user).email);
    execute format('set role %I', (auth_user).role);
end;
$$;

create or replace procedure auth.login_as_anon()
    language plpgsql
    as $$
begin
    set request.jwt.claim.sub='';
    set request.jwt.claim.role='';
    set request.jwt.claim.email='';
    set request.jwt.claims='';
    set role anon;
end;
$$;

create or replace procedure auth.logout()
    language plpgsql
    as $$
begin
    set request.jwt.claim.sub='';
    set request.jwt.claim.role='';
    set request.jwt.claim.email='';
    set request.jwt.claims='';
    set role postgres;
end;
$$;

-- pg_tap complains if there are no tests
begin;
select plan( 1 );
select ok( true, 'canary' );
select * from finish();
rollback;
