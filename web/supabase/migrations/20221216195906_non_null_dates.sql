-- Note: supabase diff produce an erroneous migration. Manually wrote this migration.
alter table user_email_sync_history alter column synced_at set not null;
alter table user_email_sync_history alter column created_at set not null;
alter table user_email_sync_history alter column updated_at set not null;

alter table user_oauth_token alter column created_at set not null;
alter table user_oauth_token alter column updated_at set not null;
alter table user_oauth_token alter column token set not null;

alter table user_profile alter column updated_at set not null;
alter table user_profile alter column created_at set not null;

alter table waitlist alter column created_at set not null;
alter table waitlist alter column updated_at set not null;
