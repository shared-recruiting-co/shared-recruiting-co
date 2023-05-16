<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `account` directory is for managing user accounts and settings. It includes components for displaying and editing user profile information, managing email settings, and rendering the account settings page. It also contains a job board component for displaying and managing job opportunities, an error page component, and a layout component for the account section of the app.

### Files
#### +error.svelte
This file contains the Svelte component for displaying an error page in the account section of the app. It imports the `PageError` component and renders it with a link to the user's profile page.

#### +layout.svelte
This file is the layout component for the account section of the web application. It imports the `Sidenav` component and renders it on the left side of the screen. The main content is rendered in a `div` element with a `slot` for dynamic content.

#### +layout.ts
This file exports a `load` function that is used to load the account layout page. The function requires the user to be logged in and checks if the user has a profile. If the user doesn't have a profile, it checks if they are on the waitlist or can create an account, and redirects them accordingly.

#### Sidenav.svelte
This file exports a Svelte component that displays a navigation sidebar for the account page. It includes links to various account-related pages, the user's profile picture and name, and a logout button. The component is responsive and includes an off-canvas menu for mobile devices.

### Directories
#### jobs
The `jobs` directory is for displaying and managing a job board in the user's account section, including a list of inbound job opportunities parsed from recruiting emails. It contains a Svelte component that allows the user to interact with the jobs by marking their interest level, applying, saving for later, or removing the job. It also defines a `load` function that retrieves a user's jobs from the database and returns them along with pagination data.

#### profile
The `profile` directory is for managing user profile data and settings. It includes a Svelte component for displaying and editing user profile information, a TypeScript file for retrieving user data, and a Svelte component for managing email settings.

#### settings
The `settings` directory is for managing the user's account settings. It includes a Svelte component for rendering the account settings page, a TypeScript module for loading the page and checking if the user is logged in, and a Svelte component for displaying a modal for deleting the user's account.

#### setup
The `setup` directory is for loading and rendering the account setup page, which includes a header, content explaining the need for connecting to a Google mail account, a checkbox for agreeing to the SRC Terms of Service, a button to connect a Google account, and links to the SRC privacy policy and code on GitHub. It also checks if the user is logged in, has a profile, and is on the waitlist, and redirects them if necessary.

<!--- END SELFDOC --->