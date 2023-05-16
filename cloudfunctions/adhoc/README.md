# Adhoc Cloud Functions

This is a catch-all package for writing cloud function workflows. For example,

- Migrating to a new label scheme
- Change label colors
- Re-running workflows on previously processed data
- Deleting @SRC data from an account
- Etc


<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `adhoc` directory is for writing cloud function workflows that perform various tasks such as migrating to a new label scheme, changing label colors, re-running workflows on previously processed data, deleting SRC data from an account, and more. It contains files for linting configuration, module definition, and helper functions for decoding JSON and filtering Gmail messages. Additionally, it includes files for specific cloud functions such as `migrateLabels`, `populateJobs`, `reclassify`, and `unsubscribe`.

### Files
#### .golangci.yml
This file contains the configuration for golangci-lint, a linter for Go code. It specifies a list of issue texts to exclude from linting, with an example of excluding a specific issue related to loop variables.

#### go.mod
This file (`cloudfunctions/adhoc/go.mod`) contains the module definition for the `adhoc` directory, including the required dependencies and their versions.

#### go.sum
This file contains a list of dependencies and their specific versions for the cloudfunctions/adhoc directory, with cloud.google.com/go being the main dependency.

#### main.go
This file contains the main function for the adhoc cloud functions. It imports necessary packages and defines the `init()` function which registers the HTTP handlers for the different cloud functions. It also includes helper functions for decoding JSON from environment variables and filtering Gmail messages.

#### migrate_labels.go
This file contains the `migrateLabels` function, which migrates old label schemes to new ones for all users with valid auth tokens. It also includes the `syncLabels` function, which updates each label to update properties.

#### populate_jobs.go
This file contains the `populateJobs` function which fetches job data from Gmail threads, filters messages before a reply, picks the earliest message, and hits the parse endpoint. If all required fields are present, it inserts the job into the `user_email_jobs` table. Otherwise, it logs the error.

#### reclassify.go
This file contains the `reclassify` function which reclassifies emails as non-recruiting by removing labels and parent folders. It also fetches recruiting threads, creates SRC labels for each user's Gmail account, and handles OAuth errors. The function uses batch processing to reduce memory usage and initializes Sentry for error tracking.

#### unsubscribe.go
This file contains the `unsubscribe` Cloud Function, which is responsible for unsubscribing users from all email notifications. It handles an HTTP request, checks for OAuth token expiration or revocation, validates the user's email is active, and uses Sentry for error tracking. If an error occurs during the process, it is logged and captured by Sentry.

<!--- END SELFDOC --->