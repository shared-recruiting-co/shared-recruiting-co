# Candidate Email Sync

`candidate_email_sync` is an HTTP triggered cloud function for syncing a user's inbox to a start and end date. It is used to do an initial one-time historic sync for new users and future syncs if a user's account is inactive for over a week.

`candidate_email_sync` blindly syncs emails to the inputted start date, which makes it convenient for manual DevOps tasks that require re-syncing a user's inbox.


<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `candidate_email_sync` directory is for a Cloud Function package that syncs candidate emails. It includes a README file with a brief overview of the function's purpose, a configuration file for golangci-lint, a file with functions that fetch threads of emails from a Gmail account and filter them based on certain criteria, a file with unit tests for those functions, a file that specifies the module name, Go version, and dependencies required for the function to run, and a file with the Cloud Function package that imports various packages, including `pubsub`, `gmail`, and `sentry-go`.

### Files
#### .golangci.yml
This file contains the configuration for golangci-lint, a tool used to run various linters on Go code. It includes a list of regex patterns to exclude from the linting process, with an example of excluding a specific issue related to loop variables being captured by a function literal.

#### cloudfunction.go
This file contains a Cloud Function package that syncs candidate emails. It imports various packages, including `pubsub`, `gmail`, and `sentry-go`. The `handler` function initializes the CloudFunction struct, sets up error handling and logging using Sentry, reads the request body, and calls the `Sync` method. The `Sync()` function fetches and processes new emails from Gmail using batch processing and waits for all messages to be processed before returning. Additionally, the file contains a function `PublishMessages` that publishes Gmail messages to a Pub/Sub topic.

#### email.go
This file contains functions that fetch threads of emails from a Gmail account and filter them based on certain criteria. The `fetchThreadsSinceDate` function fetches all threads since a given start date, ignoring threads of only sent emails and threads already processed by SRC. The `skipThread` function skips threads that have already been labeled with the SRC label, and the `filterMessagesAfterReply` function filters messages in a thread to only include those that were sent after a reply.

#### email_test.go
This file contains unit tests for the `skipThread` and `filterMessagesAfterReply` functions in the `email.go` file. The tests cover different scenarios and use the `testing` package and the `google.golang.org/api/gmail/v1` package.

#### go.mod
This file specifies the module name, Go version, and dependencies required for the `candidate_email_sync` cloud function to run, including packages such as `cloud.google.com/go/pubsub`, `github.com/GoogleCloudPlatform/functions-framework-go`, `github.com/getsentry/sentry-go`, `github.com/shared-recruiting-co/shared-recruiting-co/libs/src`, and `google.golang.org/api`.

#### go.sum
This file contains a list of dependencies and their specific versions, as well as cryptographic checksums for the modules used in the `cloudfunctions/candidate_email_sync` directory, ensuring the integrity of the downloaded modules.

<!--- END SELFDOC --->