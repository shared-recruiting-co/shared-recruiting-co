# Candidate Gmail Label Changes

`candidate_gmail_label_changes` is a handler for reacting to `@SRC` label changes (added or removed).

This cloud function serves many purposes:

- Allows users to trigger actions based on adding or removing labels in Gmail
- Allows SRC to sync state changes between the user's inbox and the user's data in the database
- Centralized logic for reacting to label change events


<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `candidate_gmail_label_changes` directory is for a Cloud Function in Go that handles changes to Gmail labels for emails related to candidates in the Shared Recruiting Co. platform. It allows users to trigger actions based on adding or removing labels in Gmail, sync state changes between the user's inbox and the user's data in the database, and provides centralized logic for reacting to label change events. The directory includes the implementation file, a README file, and go module files.

### Files
#### cloudfunction.go
This file contains the implementation of a Cloud Function in Go that handles changes to Gmail labels for emails related to candidates in the Shared Recruiting Co. platform. It imports various libraries for handling Gmail API requests, database queries, and machine learning services. The Cloud Function is triggered by Cloud Events and contains functions for handling added and removed labels. The file includes methods for parsing email, inserting data into a database, and checking if an email is a known recruiting outbound message. It also handles different Gmail label changes related to job opportunities and logs relevant information for each function. Error logging and capturing is implemented using Sentry.

#### go.mod
This file is the go module file for the `candidate_gmail_label_changes` cloud function. It lists the required dependencies and their versions.

#### go.sum
This file (`go.sum`) contains a list of dependencies and their specific versions, including cryptographic hashes, used in the `cloudfunctions/candidate_gmail_label_changes` directory and cloud function.

<!--- END SELFDOC --->