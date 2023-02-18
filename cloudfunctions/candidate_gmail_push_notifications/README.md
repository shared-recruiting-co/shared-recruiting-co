# Candidate Push Notifications

`candidate_gmail_push_notification` is a handler for [Gmail push notifications](https://developers.google.com/gmail/api/guides/push).

It is triggered every time a watched event happens in a Gmail inbox with SRC installed.

This cloud function serves many purposes:

- Keep track of most recent user inbox history ID that SRC has synced to
- Trigger a historic inbox sync if it's the user's first sync or if the history ID has expired (over one week since last sync)
- Fetch new messages since last sync 
