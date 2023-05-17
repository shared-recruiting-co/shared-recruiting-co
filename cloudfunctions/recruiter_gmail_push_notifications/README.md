# Recruiter Push Notifications

`recruiter_gmail_push_notification` is a handler for [Gmail push notifications](https://developers.google.com/gmail/api/guides/push).

It is triggered every time a watched event happens in a Gmail inbox with SRC installed.

This cloud function serves many purposes:

- Keep track of most recent user inbox history ID that SRC has synced to
- Trigger a historic inbox sync if it's the user's first sync or if the history ID has expired (over one week since last sync)
- Fetch new messages since last sync 


<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `recruiter_gmail_push_notifications` directory is for a Cloud Function that handles email push notifications in the SRC recruiting platform. It includes functions for syncing recruiter Gmail inboxes with the platform, fetching new messages since the last sync, and triggering historic inbox syncs. It also uses Sentry for error reporting.

### Files
#### .golangci.yml
This file contains the configuration for golangci-lint, a linter for Go code. It specifies an exclusion for a specific issue related to loop variables being captured by function literals.

#### cloudfunction.go
This file defines a Cloud Function for handling email push notifications in the SRC recruiting platform. It includes functions for retrieving user credentials, creating a Gmail service, and syncing recruiter Gmail inboxes with the platform. It also uses Sentry for error reporting.

#### email.go
This file contains a function `fetchSentEmailsSinceHistoryID` that fetches sent emails since a given history ID using the Gmail API. It takes in a service, history ID, and page token as arguments and returns a slice of Gmail messages, the next page token, and an error.

#### go.mod
This file contains the module dependencies for the `recruiter_gmail_push_notifications` cloud function, including required versions of various Google Cloud libraries and third-party packages.

#### go.sum
This file contains a list of dependencies and their specific versions, along with cryptographic checksums, for the `recruiter_gmail_push_notifications` cloud function. It includes dependencies such as `cloud.google.com/go/pubsub`, `cloud.google.com/go/storage`, `golang.org/x/oauth2`, and `golang.org/x/sync`. These checksums are used to ensure that the correct versions of the dependencies are installed when building and deploying the function.

#### sync.go
This file contains functions for syncing emails and publishing messages to a Pub/Sub topic in the recruiter Gmail push notifications feature. It includes a function for triggering a background sync of all emails since a given date and a function for syncing email history.

<!--- END SELFDOC --->