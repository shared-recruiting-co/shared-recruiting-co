<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `cloudfunctions` directory is for various Cloud Functions that handle different tasks related to the Shared Recruiting Co. platform. These tasks include syncing candidate and recruiter emails, handling changes to Gmail labels, processing and labeling recruiting emails, and handling Gmail push notifications. Each subdirectory contains implementation files, test functions, configuration files, and module dependencies.

### Directories
#### adhoc
The `adhoc` directory is for writing cloud function workflows that perform various tasks such as migrating to a new label scheme, changing label colors, re-running workflows on previously processed data, deleting SRC data from an account, and more. It contains files for linting configuration, module definition, and helper functions for decoding JSON and filtering Gmail messages. Additionally, it includes files for specific cloud functions such as `migrateLabels`, `populateJobs`, `reclassify`, and `unsubscribe`.

#### candidate_email_sync
The `candidate_email_sync` directory is for a Cloud Function package that syncs candidate emails. It includes a README file with a brief overview of the function's purpose, a configuration file for golangci-lint, a file with functions that fetch threads of emails from a Gmail account and filter them based on certain criteria, a file with unit tests for those functions, a file that specifies the module name, Go version, and dependencies required for the function to run, and a file with the Cloud Function package that imports various packages, including `pubsub`, `gmail`, and `sentry-go`.

#### candidate_gmail_label_changes
The `candidate_gmail_label_changes` directory is for a Cloud Function in Go that handles changes to Gmail labels for emails related to candidates in the Shared Recruiting Co. platform. It allows users to trigger actions based on adding or removing labels in Gmail, sync state changes between the user's inbox and the user's data in the database, and provides centralized logic for reacting to label change events. The directory includes the implementation file, a README file, and go module files.

#### candidate_gmail_messages
The `candidate_gmail_messages` directory is for a Google Cloud Function written in Go that processes and labels recruiting emails in a candidate's Gmail inbox. It includes functions for sorting messages by date, checking if a message thread has already been labeled, and processing and labeling new and known recruiting emails. The directory also contains the `go.mod` and `go.sum` files, which list the required dependencies and their versions, including cryptographic checksums.

#### candidate_gmail_push_notifications
The `candidate_gmail_push_notifications` directory is for a cloud function that handles Gmail push notifications. It retrieves the user's profile and OAuth token, creates a Gmail service, gets or creates SRC labels, and sets up a Pub/Sub client for sending messages. It also includes helper functions for handling errors and parsing JSON. Additionally, it contains functions for fetching changes to a Gmail account's history, converting those changes to added messages, and filtering email label changes based on label IDs.

#### gmail_subscription
The `gmail_subscription` directory is for a Cloud Function that handles Gmail subscriptions for candidates and recruiters. It includes functions for checking email activity, creating Gmail subscriptions, and error tracking with Sentry. The directory also contains configuration files for GolangCI, a README for the `watch_emails` function, and `go.mod` and `go.sum` files for managing dependencies.

#### recruiter_email_sync
The `recruiter_email_sync` directory is for a cloud function that syncs a user's inbox to a start and end date. It includes files for linting configuration, a README, the CloudFunction struct definition, functions for syncing emails and handling errors, and functions for fetching and skipping threads. It also includes files for module dependencies and their versions.

#### recruiter_gmail_messages
The `recruiter_gmail_messages` directory is for a Cloud Function that handles incoming Gmail messages for the recruiter. It includes the implementation of the function, a test function, a GolangCI configuration file, and the `go.mod` and `go.sum` files for the required packages.

#### recruiter_gmail_push_notifications
The `recruiter_gmail_push_notifications` directory is for a Cloud Function that handles email push notifications in the SRC recruiting platform. It includes functions for syncing recruiter Gmail inboxes with the platform, fetching new messages since the last sync, and triggering historic inbox syncs. It also uses Sentry for error reporting.

<!--- END SELFDOC --->