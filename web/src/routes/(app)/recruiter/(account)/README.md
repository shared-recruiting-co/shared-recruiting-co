<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `account` directory is for managing the recruiter's account section in the SRC platform. It includes a layout component, a responsive Sidenav component, and directories for managing job listings, outbound recruiting emails, profile information, and account settings.

### Files
#### +layout.svelte
This file is the layout component for the recruiter account page. It imports the Sidenav component and renders it alongside the main content of the page using a CSS grid layout. The main content is passed as a slot to this component.

#### Sidenav.svelte
This file exports a responsive Sidenav component for the recruiter account page in the SRC platform. It displays the recruiter's profile information, navigation links, and a logout button. The component includes a toggle button for mobile view and uses Svelte syntax to conditionally apply classes based on the current page.

### Directories
#### jobs
The `jobs` directory is for managing job listings in a recruiter's account section. It includes a Svelte component for displaying a table of jobs, a TypeScript file for retrieving job data, and directories for displaying job details, candidates, and outbound templates for a specific job, as well as creating a new job listing.

#### outbound
The `outbound` directory is for displaying and managing outbound recruiting emails in a table format with pagination controls, a dropdown menu with job options, and a hidden table cell with an icon. It includes a Svelte component file for displaying the emails and a TypeScript file for retrieving and limiting the number of templates displayed.

#### profile
The `profile` directory is for managing the recruiter's profile page, including a Svelte component for updating profile information, a function for retrieving Gmail accounts associated with a recruiter's profile, and a Svelte component for activating and deactivating Gmail integration.

#### settings
The `settings` directory is for managing the recruiter account settings in SRC. It contains a Svelte component that displays a welcome message and a notification for an onboarding call, as well as a link to contact the SRC team to schedule the call.

<!--- END SELFDOC --->