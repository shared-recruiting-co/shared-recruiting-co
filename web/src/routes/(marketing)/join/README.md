<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `join` directory is for handling user signups and waitlist submissions. It includes a Svelte component for the signup form, a server-side script for form validation and submission, and a function for checking if the user is already on the waitlist.

### Files
#### +page.server.ts
This file (`web/src/routes/(marketing)/join/+page.server.ts`) exports an object `actions` containing a `default` function that handles a form submission. It requires the user to be logged in and validates the form fields. If the form is valid, it adds the user to the waitlist by inserting their information into the `waitlist` table in the Supabase database. If there is an error, it returns a 400 status code with an error message, otherwise it returns a success message.

#### +page.svelte
This file is a Svelte component that defines a form for users to join SRC. It includes fields for user information such as first name, last name, LinkedIn profile, and how the user heard about SRC. It also displays success and error messages upon form submission.

#### +page.ts
This file exports a `load` function that requires the user to be logged in and checks if they are already on the waitlist. If they are, it either redirects them to the profile page or shows them the success state. If they are not on the waitlist, it shows the form.

<!--- END SELFDOC --->