<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `unsubscribe` directory is for handling the POST request to unsubscribe a user's email from Gmail notifications and updating the email settings in the database.

### Files
#### +server.ts
This file contains a POST request handler that unsubscribes a user's email from Gmail notifications and updates the email settings in the database. It checks for authorization, retrieves the user's email and Google access token, stops watching for new emails, and deactivates the email for either a candidate or a recruiter. If there is an error, it throws an error response or a 500 error if there is an error while updating the settings.

<!--- END SELFDOC --->