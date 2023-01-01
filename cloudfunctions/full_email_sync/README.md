# Full Email Sync

`full_email_sync` is an HTTP triggered cloud function for syncing a user's inbox to a give date. It is used to do an initial one-time historic sync for new users and future syncs if a user's account is inactive for over a week. The logic is nearly identical to `email_push_notification` and will be refactor at a future date.

`full_email_sync` blindly syncs emails to the inputted start date, which makes it convenient for manual DevOps tasks that require re-syncing a user's inbox.
