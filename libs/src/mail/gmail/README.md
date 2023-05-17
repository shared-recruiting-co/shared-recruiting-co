# gmail

`gmail` is a wrapper Go library for working with the [Google gmail API](https://pkg.go.dev/google.golang.org/api/gmail/v1). 

It is in active development and makes no guarantees about API stability. The long-term goal is to create a standard `mail` interface we can use across all email clients, like Gmail and Outlook.


<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `gmail` directory is for working with the Google Gmail API through a Go library. It includes functions for forwarding and fetching messages, managing labels, and handling errors. The directory also contains subdirectories for managing Gmail labels and retrieving information from Gmail messages.

### Files
#### clone.go
This file contains the `CloneMessage` function, which is used to clone a message from one inbox to another. Cloning preserves the original message sender, recipient, and all other headers. The function also supports recipient anonymization.

#### errors.go
This file contains functions and default values for handling various types of errors that may be encountered when interacting with Gmail through the Google API, including OAuth2 errors, Google API errors, rate limit errors, and not found errors. Additionally, it provides a function for executing a function with retries using an exponential backoff if a rate limit error is encountered.

#### errors_test.go
This file contains multiple test functions for different error identification functions in the `gmail` package. The tests check if errors are correctly identified as OAuth2, Google API, rate limit, or `googleapi.Error` with a `http.StatusNotFound` code. Additionally, there is a test function that iterates through a list of test cases and checks if the number of calls matches the expected value and if the error matches the expected error.

#### fwd.go
This file contains structs and functions for forwarding messages in Gmail. It includes methods for retrieving headers and message content, as well as functions for creating raw forwarded messages and sendable Gmail messages.

#### gmail.go
This file implements a Go library for working with the Google Gmail API. It includes convenience methods for interacting with Gmail, such as retrieving user profile, listing and creating labels, and forwarding and fetching messages. It also includes functions for managing labels, allowing or blocking senders, and fetching or creating labels managed by SRC.

### Directories
#### label
The `label` directory is for defining and managing Gmail labels used by SRC for organizing and filtering emails. It includes candidate-specific labels for managing job opportunities and blocked emails, as well as labels for recruiters, such as `Recruiting` and `RecruitingOutbound`.

#### message
The `message` directory is for retrieving and sorting information from Gmail messages, including sender and recipient email addresses, message subject and body, message ID, and label ID. It also includes functions for checking if a message was sent by the current user and returning the time the message was created. Additionally, it contains unit tests for various functions in the `message` package to ensure they return the expected output.

<!--- END SELFDOC --->