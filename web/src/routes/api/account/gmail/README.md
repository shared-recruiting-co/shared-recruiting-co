<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `gmail` directory is for handling various functionalities related to a user's Gmail account. It includes subdirectories for connecting a user's Gmail account to their SRC account, modifying Gmail labels on a thread, subscribing to receive notifications for unread emails, and unsubscribing from email notifications.

### Directories
#### connect
The `connect` directory is for handling requests to connect a user's Gmail account to their SRC account. It includes a file that exports an async function and two request handlers that exchange a code for access and refresh tokens, validate the scope, and save the OAuth tokens. The handlers parse the JWT payload to get the user ID, verify that the email matches, and return the email.

#### labels
The `labels` directory is for modifying Gmail labels on a thread. It includes a DELETE handler that removes all "@SRC" labels and a PUT handler that can add or remove specified labels. The handlers require an email and threadId and use a refreshed Google access token to interact with the Gmail API.

#### subscribe
The `subscribe` directory is for handling the subscription of a user's Gmail account to receive notifications for unread emails. The `+server.ts` file exports a POST request handler that checks user authorization, gets the user's email and Google refresh token, and watches for new emails to update the user's profile.

#### unsubscribe
The `unsubscribe` directory is for handling the POST request to unsubscribe a user's email from Gmail notifications and updating the email settings in the database.

<!--- END SELFDOC --->