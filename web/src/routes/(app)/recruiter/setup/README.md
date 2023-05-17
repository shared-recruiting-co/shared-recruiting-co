<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `setup` directory is for handling the recruiter setup page, including server-side code for creating a new company and recruiter in the database based on user input, and a Svelte component for displaying and handling the form for confirming profile information.

### Files
#### +page.server.ts
This file contains server-side code for the recruiter setup page. It imports functions from `$lib/forms` to validate form data submitted by the user, requires the user to be logged in, and checks if the terms of service checkbox is checked. It also creates a new company and recruiter in the database based on the user's input and returns a success message if the operation is successful, or a 400 error with an error message if there are any errors.

#### +page.svelte
This file is a Svelte component that handles the setup page for recruiters. It contains a form with fields for confirming profile information, including email and first name, as well as fields for the recruiter's name, company name, company website, and how they heard about SRC. The component also includes error handling for each field and a function to handle user logout.

<!--- END SELFDOC --->