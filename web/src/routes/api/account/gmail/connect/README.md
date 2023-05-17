<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `connect` directory is for handling requests to connect a user's Gmail account to their SRC account. It includes a file that exports an async function and two request handlers that exchange a code for access and refresh tokens, validate the scope, and save the OAuth tokens. The handlers parse the JWT payload to get the user ID, verify that the email matches, and return the email.

### Files
#### +server.ts
This file exports an async function and two request handlers that handle requests to connect a user's Gmail account to their SRC account. The function exchanges a code for access and refresh tokens, validates the scope, and returns a string. The handlers parse the JWT payload to get the user ID, verify that the email matches, save the OAuth tokens, and return the email. If there is an error, an appropriate error message is thrown. The file also defines helper functions `expiryFromExpiresIn` and `parseJWTPayload`.

<!--- END SELFDOC --->