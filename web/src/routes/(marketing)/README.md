<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `(marketing)` directory is for storing files related to the marketing section of the Shared Recruiting Co. (SRC) website. It includes components and pages for the main landing page, as well as specific pages for candidates and companies. Additionally, it contains directories for handling user signups, legal pages, login, and security.

### Files
#### +error.svelte
This file contains the Svelte component for displaying an error page in the marketing section of the website. It imports the `PageError` component from `$lib/components/PageError.svelte` and renders it in a flex container that centers the content on the page.

#### +layout.svelte
This file contains the layout component for the marketing pages. It imports the Header and Footer components and renders them around the main content of the page, which is passed in as a slot.

#### +page.svelte
This file (`+page.svelte`) contains the main marketing page for the Shared Recruiting Co. (SRC) website. It includes a hero section with a call-to-action for candidates and companies to join the platform, as well as information about SRC's mission and values. Additionally, it features a section highlighting the benefits of the platform for candidates, companies, and recruiting.

#### Footer.svelte
This file contains the `Footer` component of the marketing pages for SRC, which includes links to various sections of the website, social media accounts, email, and legal documents, as well as displaying copyright information.

#### Header.svelte
This file contains the `Header` component for the marketing pages. It is responsible for rendering the navigation bar at the top of the page, including a logo and navigation links for candidates, companies, and security. It also includes conditional rendering for the account and logout buttons based on whether the user is logged in or not, a "Join" button that links to the login page, and a GitHub icon that links to the SRC GitHub repository. Additionally, it includes a mobile menu button and a function to handle user logout.

### Directories
#### candidates
The `candidates` directory is for Svelte components and pages related to the main landing page of the Shared Recruiting Co. (SRC) candidate-centric recruiting platform. It includes a Svelte page component for the landing page, a component for displaying features of the platform, and a component for displaying a testimonial from a candidate.

#### companies
The `companies` directory is for files related to the marketing page of the Shared Recruiting Co. (SRC) website that are specific to companies. It includes the main landing page (`+page.svelte`), a component for displaying features (`Features.svelte`), and a component for displaying statistics (`Stats.svelte`).

#### join
The `join` directory is for handling user signups and waitlist submissions. It includes a Svelte component for the signup form, a server-side script for form validation and submission, and a function for checking if the user is already on the waitlist.

#### legal
The `legal` directory is for storing legal pages and their content for the SRC recruiting platform. It includes subdirectories for `privacy-policy` and `terms-of-service`, each with a Svelte component that renders the page and a Markdown file that contains the corresponding legal content.

#### login
The `login` directory is for handling the login page of the marketing section of the SRC platform. It includes a Svelte component for the login page with a form for users to sign in with their Google account and an error message display. Additionally, it contains a TypeScript file with a `load` function that checks if the user is logged in and redirects them to the appropriate page.

#### security
The `security` directory is for storing the markup and content of the security page on the Shared Recruiting Co. marketing website. It includes a hero section highlighting the platform's privacy-first approach, as well as sections discussing SRC's commitment to protecting user data, compliance with the Cloud Application Security Assessment (CASA), and the importance of public scrutiny in identifying and patching security issues.

<!--- END SELFDOC --->