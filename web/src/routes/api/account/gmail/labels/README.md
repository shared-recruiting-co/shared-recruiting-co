<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `labels` directory is for modifying Gmail labels on a thread. It includes a DELETE handler that removes all "@SRC" labels and a PUT handler that can add or remove specified labels. The handlers require an email and threadId and use a refreshed Google access token to interact with the Gmail API.

### Files
#### +server.ts
This file exports request handlers for modifying Gmail labels on a thread. It includes a DELETE handler that removes all "@SRC" labels and a PUT handler that can add or remove specified labels. The handlers require an email and threadId and use a refreshed Google access token to interact with the Gmail API.

<!--- END SELFDOC --->