# Candidate Gmail Label Changes

`candidate_gmail_label_changes` is a handler for reacting to `@SRC` label changes (added or removed).

This cloud function serves many purposes:

- Allows users to trigger actions based on adding or removing labels in Gmail
- Allows SRC to sync state changes between the user's inbox and the user's data in the database
- Centralized logic for reacting to label change events
