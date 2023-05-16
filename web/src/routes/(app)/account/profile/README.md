<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `profile` directory is for managing user profile data and settings. It includes a Svelte component for displaying and editing user profile information, a TypeScript file for retrieving user data, and a Svelte component for managing email settings.

### Files
#### +page.svelte
This file (`+page.svelte`) contains a Svelte component that handles user profile data, displays account information, inbound and verified jobs, and a form for editing the user's profile. It also includes functions for debouncing user input and formatting dates.

#### +page.ts
This file exports a `load` function that retrieves and returns data related to the user's email sync history, OAuth token validity, and job counts. It requires the user to be logged in and throws a redirect error if not. The returned data includes the last time the user synced their emails, whether their OAuth token is valid, and the number of inbound and official jobs.

#### EmailSettings.svelte
This file contains the Svelte component for the email settings page. It imports various components and functions, including `Toggle`, `Hint`, `AlertModal`, and `ConnectGoogleAccountButton`. It exports `isValid`, `email`, `settings`, and `saveSettings` as props, and contains various UI states and functions for toggling and handling deactivation confirmation.

<!--- END SELFDOC --->