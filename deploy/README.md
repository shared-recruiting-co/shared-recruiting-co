<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `deploy` directory is for deploying the Shared Recruiting Co. platform to Google Cloud Platform using Pulumi and Go code. It includes configuration files for Pulumi and GolangCI, Go code for creating and deploying various Google Cloud Functions, and files for managing dependencies and testing.

### Files
#### .golangci.yml
This file contains the configuration for GolangCI, a linter for Go code. It includes a list of regex patterns to exclude from linting, with an example of excluding a loop variable captured by a function literal.

#### Pulumi.prod.yaml
This file contains the Pulumi stack configuration in YAML format for deploying the Shared Recruiting Co. platform to Google Cloud Platform, including secure credentials for various services.

#### Pulumi.yaml
This file contains the configuration for Pulumi, a tool for deploying cloud infrastructure. It specifies the name of the project, the runtime language, and a brief description of the infrastructure.

#### cloudfunction.go
This file contains Go code for creating and deploying various Google Cloud Functions, including functions for syncing candidate and recruiter Gmail accounts, handling label changes, and subscribing to email inboxes. The code sets up the functions' build and service configurations, environment variables, and access permissions.

#### cloudfunction_test.go
This file contains a test function for the `shortenAccountId` function in the `cloudfunction.go` file. It tests the function's ability to shorten a given service account name to a specific length and replace common names. The test uses a table-driven approach to test multiple inputs and expected outputs.

#### go.mod
This file lists the dependencies and their versions for the project, including the Pulumi SDK and the Pulumi GCP SDK.

#### go.sum
This file (`deploy/go.sum`) contains a list of dependencies and their specific versions, including cryptographic hashes, used in the project's deployment and dependencies. It ensures that the correct versions of dependencies are installed when building the project.

#### main.go
This file contains the `main` function for deploying the Shared Recruiting Co. platform. It sets up the infrastructure for the project, including creating a storage bucket, setting up topics, and creating cloud functions.

#### pubsub.go
This file creates a protected Pub/Sub topic for candidate Gmail label changes with enforced schema validation using the Pulumi GCP SDK. It returns an error if the topic creation fails.

<!--- END SELFDOC --->