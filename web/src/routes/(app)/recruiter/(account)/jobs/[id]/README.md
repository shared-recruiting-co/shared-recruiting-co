<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `[id]` directory is for displaying job details, candidates, and outbound templates for a specific job in a recruiter's account. It includes a Svelte component for the job details page, a TypeScript file for fetching job data, a Svelte component for displaying a list of candidates, and a Svelte component for displaying outbound templates assigned to the job.

### Files
#### +page.svelte
This file is a Svelte component that displays the job details page for a recruiter account. It includes navigation tabs for candidates and outbound templates and renders the corresponding components based on the URL hash. Additionally, it contains TODOs for editing and deleting the job.

#### +page.ts
This file exports a function that fetches job data from the database, including information about candidates and outbound templates. It also maps over the outbound templates list to add a recipient count property to each one. If there is an issue fetching the data, an error is logged.

#### Candidates.svelte
This file (`Candidates.svelte`) contains a Svelte component that displays a list of candidates for a specific job. It renders them in a table and displays an empty state with a message and an icon if there are no candidates. The table shows the candidate's email and last activity date. It also has a TODO comment to determine which columns to show in the table.

#### OutboundTemplates.svelte
This file, `OutboundTemplates.svelte`, defines a Svelte component responsible for displaying outbound templates assigned to a job posting in the recruiter's account. The component exports an array of outbound templates with their corresponding recipient count and renders them in a table with columns for subject, body, recipient count, and last sent date. If there are no templates assigned, an empty state is displayed with a message prompting the user to assign outbound templates.

<!--- END SELFDOC --->