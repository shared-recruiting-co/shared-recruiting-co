<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `google` directory is for interacting with Google APIs, specifically the Gmail and OAuth2 APIs. It includes functions for sending and receiving emails, exchanging authorization codes for access and refresh tokens, and refreshing access tokens. Additionally, it includes a function for sending a welcome email to new or returning users of Shared Recruiting Co.

### Files
#### gmail.ts
This file contains functions and helper functions for interacting with Gmail using the Gmail API. It includes functions to watch for new emails, send email messages, retrieve label IDs, and check if a thread ID corresponds to a Gmail thread. It also includes helper functions to retrieve and update label IDs of a Gmail thread.

#### oauth.ts
This file contains functions for exchanging an authorization code for access and refresh tokens, and refreshing an access token using the Google OAuth2 API.

#### welcomeEmail.ts
This file exports a function `sendWelcomeEmail` that uses the Google API to send a welcome email to a new or returning user of Shared Recruiting Co. The function refreshes the user's access token and constructs the email message. If there is an error, it fails silently.

<!--- END SELFDOC --->