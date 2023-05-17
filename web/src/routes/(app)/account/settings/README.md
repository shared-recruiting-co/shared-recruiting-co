<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `settings` directory is for managing the user's account settings. It includes a Svelte component for rendering the account settings page, a TypeScript module for loading the page and checking if the user is logged in, and a Svelte component for displaying a modal for deleting the user's account.

### Files
#### +page.svelte
This file is a Svelte component that renders the account settings page. It includes a "Delete Account" button that, when clicked, opens a modal component (DeleteAccountModal.svelte) that allows the user to confirm the deletion of their account.

#### +page.ts
This file contains a TypeScript module that exports a `load` function, which is a SvelteKit `PageLoad` function. The `load` function requires the user to be logged in and throws a redirect to the login page if the user is not logged in. If the user is logged in, an empty object is returned.

#### DeleteAccountModal.svelte
This file contains the `DeleteAccountModal` Svelte component which displays a modal for deleting a user's account. It includes a form for entering a reason for deletion and handles the deletion process via an API call. The modal is displayed using the `show` boolean prop and includes a fade-in and fade-out transition. The component also includes conditional rendering for an error message and a loading state.

<!--- END SELFDOC --->