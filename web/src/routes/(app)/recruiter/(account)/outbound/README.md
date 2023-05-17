<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `outbound` directory is for displaying and managing outbound recruiting emails in a table format with pagination controls, a dropdown menu with job options, and a hidden table cell with an icon. It includes a Svelte component file for displaying the emails and a TypeScript file for retrieving and limiting the number of templates displayed.

### Files
#### +page.svelte
This file is a Svelte component that displays outbound recruiting emails in a table with columns for subject, body, and job. It includes pagination controls, a dropdown menu with job options, and a hidden table cell with an icon. If no emails are found, a message prompts the user to connect their Gmail account. The file also defines a function `onSelect` that updates the `job_id` field in the `recruiter_outbound_template` table.

#### +page.ts
This file exports a `load` function that retrieves outbound recruiter templates from the database, limits the number of templates displayed using pagination, and updates the sequences with recipient counts. It also includes a TODO comment suggesting potential client-side fetching of jobs in the combo box.

<!--- END SELFDOC --->