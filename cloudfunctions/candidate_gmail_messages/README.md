<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `candidate_gmail_messages` directory is for a Google Cloud Function written in Go that processes and labels recruiting emails in a candidate's Gmail inbox. It includes functions for sorting messages by date, checking if a message thread has already been labeled, and processing and labeling new and known recruiting emails. The directory also contains the `go.mod` and `go.sum` files, which list the required dependencies and their versions, including cryptographic checksums.

### Files
#### cloudfunction.go
This file contains Go code for a Google Cloud Function that processes and labels recruiting emails in a candidate's Gmail inbox. It includes functions for checking if a message thread has already been labeled, sorting messages by date, and processing and labeling new and known recruiting emails. It also imports various libraries for handling Gmail messages, Cloud Events, and error handling.

#### go.mod
This file is the go module file for the `candidate_gmail_messages` cloud function. It lists all the required dependencies and their versions.

#### go.sum
This file (`go.sum`) contains a list of dependencies and their specific versions, including cryptographic checksums, for the `candidate_gmail_messages` cloud function. These dependencies include modules such as `cloud.google.com/go`, `golang.org/x/crypto`, `golang.org/x/exp`, and `golang.org/x/image`. The checksums are used to ensure the integrity of the dependencies and prevent tampering.

<!--- END SELFDOC --->