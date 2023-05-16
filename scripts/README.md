<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `scripts` directory is for storing scripts related to the development and maintenance of the SRC codebase. This includes a script called `update-go-mods.sh` which updates the dependencies of several Go packages in the `cloudfunctions` directory and runs tests.

### Files
#### update-go-mods.sh
This file updates the dependencies of several Go packages in the `cloudfunctions` directory and runs tests. It uses `go mod tidy` to update the dependencies and `go get` to fetch the latest version of the `libs/src` package from the `origin/main` branch.

<!--- END SELFDOC --->