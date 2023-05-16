<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `ml` directory is for machine learning functionality related to recruiting emails and job postings. It includes a package with a client struct for classifying and parsing these documents, as well as a test file with unit tests for various methods in the package.

### Files
#### ml.go
This file defines the ML package, which includes a client struct with methods for classifying and parsing recruiting emails and job postings. The client struct also includes a NewService function for creating a new client and uses HTTP requests to communicate with the SRC API.

#### ml_test.go
This file contains unit tests for various methods in the `ml` package, including `ServiceClassify`, `BatchClassify`, and `ParseJob`. The tests create mock servers and verify that the functions correctly send HTTP requests and return expected responses.

<!--- END SELFDOC --->