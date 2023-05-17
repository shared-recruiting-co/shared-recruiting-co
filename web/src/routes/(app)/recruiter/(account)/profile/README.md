<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `profile` directory is for managing the recruiter's profile page, including a Svelte component for updating profile information, a function for retrieving Gmail accounts associated with a recruiter's profile, and a Svelte component for activating and deactivating Gmail integration.

### Files
#### +page.svelte
This file contains a Svelte component for the recruiter's profile page. It includes a form for updating the recruiter's profile information, a section for email integration with Gmail, and various features such as syncing and importing candidates, creating and managing Gmail labels, and detecting and matching email sequences to open roles. It also includes error handling and success messages.

#### +page.ts
This file exports a `load` function that retrieves a list of Gmail accounts associated with a recruiter's profile from the Supabase database. The function returns an object containing the retrieved `gmailAccounts` data.

#### GmailIntegration.svelte
This file contains a Svelte component for GmailIntegration, which allows recruiters to activate and deactivate Gmail integration for their account. It includes functions for toggling visibility, handling activation and deactivation, and displaying error messages.

<!--- END SELFDOC --->