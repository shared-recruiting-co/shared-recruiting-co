<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `jobs` directory is for managing job listings in a recruiter's account section. It includes a Svelte component for displaying a table of jobs, a TypeScript file for retrieving job data, and directories for displaying job details, candidates, and outbound templates for a specific job, as well as creating a new job listing.

### Files
#### +page.svelte
This file defines the Svelte component for the Jobs page in the recruiter's account section. It displays a table of jobs with their titles, candidate counts, and last updated dates, and includes a link to view each job and a pagination component at the bottom. Additionally, it provides a button to create a new job if there are no jobs in the account.

#### +page.ts
This file defines a `load` function that retrieves job data from the database and returns it along with pagination information. The function requires the user to be logged in and handles errors for job and candidate count retrieval. The retrieved job data includes a count of the number of candidates associated with each job.

### Directories
#### [id]
The `[id]` directory is for displaying job details, candidates, and outbound templates for a specific job in a recruiter's account. It includes a Svelte component for the job details page, a TypeScript file for fetching job data, a Svelte component for displaying a list of candidates, and a Svelte component for displaying outbound templates assigned to the job.

#### new
The `new` directory is for creating a new job listing on a company's job board. It contains a server-side file for validating and adding the job to the database, as well as a Svelte component for rendering a form with fields for job title and description URL.

<!--- END SELFDOC --->