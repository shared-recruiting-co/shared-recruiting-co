<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `src` directory is for managing the configuration and source code of the Shared Recruiting Co. platform. It includes configuration files for GolangCI and SQLC, a Go module definition file, and a file for managing dependencies. It also contains packages for interacting with the SRC database, Gmail API, and machine learning functionality related to recruiting emails and job postings. Additionally, it includes a directory for implementing the pub/sub system in the SRC platform.

### Files
#### .golangci.yml
This file contains the configuration for GolangCI, a linter for Go code. It includes a list of regular expressions to exclude certain issue texts from being flagged as errors.

#### go.mod
This file (`go.mod`) contains the module definition for the `libs/src` directory of the Shared Recruiting Co. platform. It specifies the required dependencies for the module and their respective versions.

#### go.sum
This file (`go.sum`) contains a list of all the dependencies used in the SRC project along with their specific versions and cryptographic hashes. It is used to ensure the integrity of the downloaded dependencies and prevent unauthorized changes or tampering.

#### sqlc.yaml
This file (`libs/src/sqlc.yaml`) contains the configuration for SQLC, a tool that generates type-safe Go code from SQL. It specifies the database engine, the location of the SQL queries and schema files, and the Go package and output directory for the generated code. It also includes overrides for certain database types to use the `guregu/null` package for null types.

#### src.go
This file is a Go source code file that defines the package name as "src".

### Directories
#### db
The `db` directory is for a Go database client library for SRC. It includes a package for interacting with a SQL database, functions for making HTTP requests to interact with the PostgREST API, several structs used to represent data in the SRC database, SQL queries for retrieving and updating user and OAuth token information, and SQL code to create and modify tables, policies, triggers, and functions related to various entities.

#### mail
The `mail` directory is for working with email functionality in the Shared Recruiting Co. platform. It includes a subdirectory `gmail` for interacting with the Google Gmail API through a Go library. The `gmail` directory provides functions for forwarding and fetching messages, managing labels, and handling errors. It also contains subdirectories for managing Gmail labels and retrieving information from Gmail messages.

#### ml
The `ml` directory is for machine learning functionality related to recruiting emails and job postings. It includes a package with a client struct for classifying and parsing these documents, as well as a test file with unit tests for various methods in the package.

#### pubsub
The `pubsub` directory is for implementing the pub/sub system in the SRC platform. It contains a `schema` directory for defining structs and constants used as payloads for different types of events in the pub/sub system.

<!--- END SELFDOC --->