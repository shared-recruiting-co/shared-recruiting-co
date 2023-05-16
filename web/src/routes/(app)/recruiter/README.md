<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `recruiter` directory is for managing the recruiter section of the SRC platform. It includes directories for managing the recruiter's account, login, and setup pages. It also contains a layout component for loading the recruiter's profile and company information.

### Files
#### +layout.ts
This file contains the `load` function that is used to load the recruiter layout. It checks if the user is logged in and redirects to the login page if not. It also retrieves the recruiter's profile and company information and redirects to the profile page if the current path is not an account page.

### Directories
#### (account)
The `account` directory is for managing the recruiter's account section in the SRC platform. It includes a layout component, a responsive Sidenav component, and directories for managing job listings, outbound recruiting emails, profile information, and account settings.

#### login
The `login` directory is for a Svelte component that renders a recruiter login page. It includes a sign-in form, a button to sign in with Google, and a link to the general login page. The component also includes a script to handle the sign-in process and display any errors.

#### setup
The `setup` directory is for handling the recruiter setup page, including server-side code for creating a new company and recruiter in the database based on user input, and a Svelte component for displaying and handling the form for confirming profile information.

<!--- END SELFDOC --->