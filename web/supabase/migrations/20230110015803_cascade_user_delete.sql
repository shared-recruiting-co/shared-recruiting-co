-- This script was manually created because the automatic diff didn't pick up fkey changes
 alter table public.waitlist
drop constraint waitlist_user_id_fkey,
add constraint waitlist_user_id_fkey foreign key (user_id) references auth.users(id) on delete cascade;

alter table public.user_email_stat
drop constraint user_email_stat_user_id_fkey,
add constraint user_email_stat_user_id_fkey foreign key (user_id) references auth.users(id) on delete cascade;

alter table public.user_email_sync_history
drop constraint user_email_sync_history_user_id_fkey,
add constraint user_email_sync_history_user_id_fkey foreign key (user_id) references auth.users(id) on delete cascade;

alter table public.user_profile
drop constraint user_profile_user_id_fkey,
add constraint user_profile_user_id_fkey foreign key (user_id) references auth.users(id) on delete cascade;

alter table public.user_oauth_token
drop constraint user_oauth_token_user_id_fkey,
add constraint user_oauth_token_user_id_fkey foreign key (user_id) references auth.users(id) on delete cascade;
