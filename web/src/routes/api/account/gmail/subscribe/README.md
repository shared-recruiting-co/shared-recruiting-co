<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `subscribe` directory is for handling the subscription of a user's Gmail account to receive notifications for unread emails. The `+server.ts` file exports a POST request handler that checks user authorization, gets the user's email and Google refresh token, and watches for new emails to update the user's profile.

### Files
#### +server.ts
This file exports a POST request handler that subscribes a user's Gmail account to receive notifications for unread emails. It first checks if the user is authorized, then gets the user's email and Google refresh token, and finally watches for new emails and updates the user's profile accordingly.

<!--- END SELFDOC --->