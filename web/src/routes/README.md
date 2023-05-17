<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `routes` directory is for managing the routes and pages of the SRC platform. It includes files and subdirectories for displaying error pages, managing server-side data, and creating the layout for all pages. It also includes directories for managing user accounts and settings, the marketing section of the website, API functionalities, and documentation pages.

### Files
#### +error.svelte
This file contains the Svelte component that displays an error page when there is an error in the application. It imports the `PageError` component from `$lib/components/PageError.svelte` and renders it in the center of the screen.

#### +layout.server.ts
This file exports a function that loads the server-side data for the layout component. It retrieves the session data using the `getSession` function from the `locals` object passed as an argument. The retrieved session data is returned as an object.

#### +layout.svelte
This file exports a Svelte component that serves as the layout for all pages in the application. It imports and uses various modules such as `$app/navigation`, `@vercel/analytics`, and `svelte`. It also contains a console message for developers and an event listener for authentication state changes.

#### +layout.ts
This file exports a `load` function that creates a Supabase client and retrieves the user's session. It depends on the `supabase:auth` dependency and takes in `fetch`, `data`, and `depends` parameters.

### Directories
#### (app)
The `(app)` directory is for managing the main application routes of the SRC platform. It includes subdirectories for managing user accounts and settings, as well as the recruiter section of the platform. These directories contain components for displaying and editing user and recruiter information, managing email settings, and rendering the account and setup pages. They also include layout components for loading user and recruiter profiles and company information.

#### (marketing)
The `(marketing)` directory is for storing files related to the marketing section of the Shared Recruiting Co. (SRC) website. It includes components and pages for the main landing page, as well as specific pages for candidates and companies. Additionally, it contains directories for handling user signups, legal pages, login, and security.

#### api
The `api` directory is for handling various functionalities related to user accounts and candidate profiles. The `account` directory includes functions for deleting a user's account and handling Gmail account-related functionalities. The `candidate` directory includes functions for creating user profiles and sending welcome emails to new candidates.

#### docs
The `docs` directory is for the documentation pages related to SRC's recruiting platform. It includes Svelte components for displaying the layout, header, navigation, and error pages. It also contains TypeScript modules and functions for loading and navigating the documentation pages. The directory includes a subdirectory for dynamically loading markdown files and rendering the documentation pages.

<!--- END SELFDOC --->