<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `jobs` directory is for displaying and managing a job board in the user's account section, including a list of inbound job opportunities parsed from recruiting emails. It contains a Svelte component that allows the user to interact with the jobs by marking their interest level, applying, saving for later, or removing the job. It also defines a `load` function that retrieves a user's jobs from the database and returns them along with pagination data.

### Files
#### +page.server.ts
This file defines a `load` function that retrieves a user's jobs from the database and returns them along with pagination data. It requires the user to be logged in and uses SupabaseAdmin as a workaround to allow candidates to see recruiter's name, email, and company info. The function also throws errors if there are any issues with retrieving the data.

#### +page.svelte
This file is a Svelte component that displays a job board in the user's account section, including a list of inbound job opportunities parsed from recruiting emails. It allows the user to interact with the jobs by marking their interest level, applying, saving for later, or removing the job.

<!--- END SELFDOC --->