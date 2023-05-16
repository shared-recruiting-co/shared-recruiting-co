<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `candidate` directory is for handling the creation of user profiles and sending welcome emails to new candidates. It includes a server file that verifies authorization, terms of service agreement, waitlist status, and email connection before creating the profile. The welcome email function subscribes the candidate to Gmail notifications and logs any errors.

### Files
#### +server.ts
This file contains handlers for creating a user profile and sending a welcome email to new candidates. The profile creation handler verifies authorization, terms of service agreement, waitlist status, and email connection before creating the profile. The welcome email function subscribes the candidate to Gmail notifications and logs any errors. There is a TODO to replace the homebrew email solution with a transactional email service.

<!--- END SELFDOC --->