# Candidate Email Sync

`candidate_email_sync` is an HTTP triggered cloud function for syncing a user's inbox to a start and end date. It is used to do an initial one-time historic sync for new users and future syncs if a user's account is inactive for over a week.

`candidate_email_sync` blindly syncs emails to the inputted start date, which makes it convenient for manual DevOps tasks that require re-syncing a user's inbox.


<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `recruiter_email_sync` directory is for a cloud function that syncs a user's inbox to a start and end date. It includes files for linting configuration, a README, the CloudFunction struct definition, functions for syncing emails and handling errors, and functions for fetching and skipping threads. It also includes files for module dependencies and their versions.

### Files
#### .golangci.yml
This file contains the configuration for GolangCI, a linter for Go code. It specifies a list of issue texts to exclude from linting, with an example exclusion provided.

#### cloudfunction.go
This file defines a CloudFunction struct and contains functions for syncing emails between Gmail and Pub/Sub, publishing email messages to a Pub/Sub topic, and handling errors. It also initializes Sentry for error tracking and sets up error handling for reading and unmarshalling the request body.

#### email.go
This file contains two functions: `fetchThreadsSinceDate` and `skipThread`. `fetchThreadsSinceDate` fetches all threads since a given start date, ignoring threads of only sent emails and threads already processed by SRC. `skipThread` checks if a thread has already been labeled with a given label.

#### email_test.go
This file contains a test function for the `skipThread` function in the `email.go` file. It tests three scenarios: when the thread should not be skipped, when the thread is sent from the user, and when the thread is already synced. The test uses the Google Gmail API and runs through a series of tests to ensure the `skipThread` function works as expected.

#### go.mod
This file lists the module dependencies and their version numbers for the `recruiter_email_sync` cloud function.

#### go.sum
This file (`go.sum`) contains a list of all the dependencies and their versions used in the `cloudfunctions/recruiter_email_sync` directory, along with their cryptographic checksums to ensure their integrity.

<!--- END SELFDOC --->