ALTER TABLE IF EXISTS user_email_sync_history 
    ADD COLUMN email text;

CREATE TYPE inbox_type AS ENUM ('candidate', 'recruiter');

ALTER TABLE IF EXISTS user_email_sync_history 
    ADD COLUMN inbox_type inbox_type NOT NULL DEFAULT 'candidate';
