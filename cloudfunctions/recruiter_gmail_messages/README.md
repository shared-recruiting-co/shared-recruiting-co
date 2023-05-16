<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `recruiter_gmail_messages` directory is for a Cloud Function that handles incoming Gmail messages for the recruiter. It includes the implementation of the function, a test function, a GolangCI configuration file, and the `go.mod` and `go.sum` files for the required packages.

### Files
#### .golangci.yml
This file contains the configuration for GolangCI, a linter for Go code. It specifies a list of issue texts to exclude from linting, with an example exclusion provided.

#### cloudfunction.go
This file contains the implementation of a Cloud Function that handles incoming Gmail messages for the recruiter. It imports various packages and libraries, including Google Cloud Platform's Gmail API, Sentry, and uuid. The function processes email messages received from Pub/Sub, validates them, and processes each message by checking if it matches a predefined template or is a recruiting email. It saves the message to the database and labels the thread if it matches a template or is a recruiting email.

#### cloudfunction_test.go
This file contains a test function for the `normalizeBody` function in the `cloudfunction.go` file. The test function includes multiple test cases to ensure that the `normalizeBody` function correctly normalizes the HTML body of an email.

#### go.mod
This file is the `go.mod` file for the `recruiter_gmail_messages` cloud function. It specifies the module name, version of Go, and the required packages and their versions needed to build and run the function.

#### go.sum
This file (`go.sum`) contains a list of cryptographic checksums for the Go modules used in the `cloudfunctions/recruiter_gmail_messages` directory, ensuring their integrity and expected versions.

<!--- END SELFDOC --->