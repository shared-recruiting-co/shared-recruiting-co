<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `account` directory is for handling user account-related functionalities. It includes a `DELETE` function to delete a user's account, a function to send an email to the user's email address when their account is deleted, and subdirectories for handling various functionalities related to a user's Gmail account.

### Files
#### +server.ts
This file exports a `DELETE` function that handles a DELETE request to delete a user's account. The function validates the request body, unsubscribes the user from email notifications, deletes the user's account, and sends a confirmation email to the user. If any errors occur, they are logged and captured by Sentry.

#### delete.ts
This file exports a function that sends an email to the user's email address when their SRC account is deleted. The email includes the user's email and the reason for the deletion. The function also sanitizes user input.

### Directories
#### gmail
The `gmail` directory is for handling various functionalities related to a user's Gmail account. It includes subdirectories for connecting a user's Gmail account to their SRC account, modifying Gmail labels on a thread, subscribing to receive notifications for unread emails, and unsubscribing from email notifications.

<!--- END SELFDOC --->