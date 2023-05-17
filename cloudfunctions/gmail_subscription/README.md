# Watch Emails

`watch_emails` [watches](https://developers.google.com/gmail/api/reference/rest/v1/users/watch) for new messages in a user's inbox. As advised by Google, it is triggered daily.


<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `gmail_subscription` directory is for a Cloud Function that handles Gmail subscriptions for candidates and recruiters. It includes functions for checking email activity, creating Gmail subscriptions, and error tracking with Sentry. The directory also contains configuration files for GolangCI, a README for the `watch_emails` function, and `go.mod` and `go.sum` files for managing dependencies.

### Files
#### .golangci.yml
This file contains the configuration for GolangCI, a linter for Go code. It specifies a list of issue texts to exclude from linting, with an example exclusion for a specific loop variable.

#### cloudfunction.go
This file contains Go code for a Cloud Function that handles Gmail subscriptions for candidates and recruiters. It includes functions for checking email activity, creating Gmail subscriptions, and error tracking with Sentry.

#### go.mod
This file is the `go.mod` file for the `cloudfunctions/gmail_subscription` directory. It specifies the module name and version, as well as the required dependencies and their versions.

#### go.sum
This file (`go.sum`) contains a list of all the dependencies and their versions used in the `cloudfunctions/gmail_subscription` directory, along with their cryptographic hashes. It is used to ensure the integrity of the dependencies and prevent tampering.

<!--- END SELFDOC --->