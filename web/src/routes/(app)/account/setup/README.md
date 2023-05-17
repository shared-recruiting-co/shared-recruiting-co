<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `setup` directory is for loading and rendering the account setup page, which includes a header, content explaining the need for connecting to a Google mail account, a checkbox for agreeing to the SRC Terms of Service, a button to connect a Google account, and links to the SRC privacy policy and code on GitHub. It also checks if the user is logged in, has a profile, and is on the waitlist, and redirects them if necessary.

### Files
#### +page.ts
This file contains the `load` function that is responsible for loading the account setup page. It checks if the user is logged in, if they have a profile, and if they are on the waitlist and can create an account. If any of these conditions are not met, the function will redirect the user to the appropriate page. If all conditions are met, it will return an object with the user's first name, last name, and email.

#### +page@.svelte
This file is a Svelte component that renders the account setup page, including a header with a logo and logout button, content explaining the need for connecting to a Google mail account, a checkbox for agreeing to the SRC Terms of Service, a button to connect a Google account, and links to the SRC privacy policy and code on GitHub. It also provides information about how SRC will label job opportunities in the user's inbox.

<!--- END SELFDOC --->