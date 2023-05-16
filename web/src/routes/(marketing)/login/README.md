<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `login` directory is for handling the login page of the marketing section of the SRC platform. It includes a Svelte component for the login page with a form for users to sign in with their Google account and an error message display. Additionally, it contains a TypeScript file with a `load` function that checks if the user is logged in and redirects them to the appropriate page.

### Files
#### +page.svelte
This file is the Svelte component for the login page of the marketing section of the SRC platform. It contains a form for users to sign in with their Google account and displays an error message if there is an issue with the sign-in process. It also includes a link for recruiters to sign in.

#### +page.ts
This file contains a `load` function that is executed when the login page is loaded. It checks if the user is logged in and if they have a user profile. If the user is not logged in or has a profile, they are redirected to the appropriate page. If the user is on the waitlist and can create an account, they are redirected to the account setup page, otherwise they are redirected to the waitlist page.

<!--- END SELFDOC --->