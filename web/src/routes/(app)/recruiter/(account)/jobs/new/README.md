<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `new` directory is for creating a new job listing on a company's job board. It contains a server-side file for validating and adding the job to the database, as well as a Svelte component for rendering a form with fields for job title and description URL.

### Files
#### +page.server.ts
This file contains server-side actions for creating a new job listing. It requires the user to be logged in and validates the job title and description URL. If the form is valid, it adds the job to the database and redirects the user to the jobs page.

#### +page.svelte
This file is a Svelte component that exports a form object for adding a new job to a company's job board. It contains fields for job title and description URL, as well as error handling for form submission.

<!--- END SELFDOC --->