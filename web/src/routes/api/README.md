<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `api` directory is for handling various functionalities related to user accounts and candidate profiles. The `account` directory includes functions for deleting a user's account and handling Gmail account-related functionalities. The `candidate` directory includes functions for creating user profiles and sending welcome emails to new candidates.

### Directories
#### account
The `account` directory is for handling user account-related functionalities. It includes a `DELETE` function to delete a user's account, a function to send an email to the user's email address when their account is deleted, and subdirectories for handling various functionalities related to a user's Gmail account.

#### candidate
The `candidate` directory is for handling the creation of user profiles and sending welcome emails to new candidates. It includes a server file that verifies authorization, terms of service agreement, waitlist status, and email connection before creating the profile. The welcome email function subscribes the candidate to Gmail notifications and logs any errors.

<!--- END SELFDOC --->