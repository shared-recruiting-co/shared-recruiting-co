<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `lib` directory is for exporting utility functions related to forms and pagination, Svelte components, Google APIs and services, server-side logic and APIs, and Supabase interfaces and composite types.

### Files
#### forms.ts
This file contains utility functions related to forms. It includes a debounce function, a function to retrieve form errors, a function to check if a string is a valid URL, a function to retrieve a trimmed form value, and a function to retrieve a boolean value from a checkbox input.

#### pagination.ts
This file contains helper functions for pagination. It includes a function to construct a pagination object based on the current page, results per page, and total number of results. It also exports a function to generate an array of strings representing the page numbers to display in a pagination component.

### Directories
#### components
The `components` directory is for exporting Svelte components that are used in the SRC platform. These components include a modal, a Google account connection button, a hint icon with a tooltip, a page error message, a table footer with pagination controls, and a customizable toggle switch. It also includes directories for rendering markdown and marketing-related content.

#### google
The `google` directory is for storing files related to Google APIs and services. Specifically, the `labels.ts` file defines an enum and type for labels used in the `@SRC` namespace for jobs.

#### server
The `server` directory is for handling server-side logic and APIs for the Shared Recruiting Co. platform. It includes subdirectories for interacting with external services such as Google APIs and for handling user authentication and authorization.

#### supabase
The `supabase` directory is for TypeScript interfaces, enums, and composite types used in the Supabase library, including the structure of tables, views, and functions in the database, as well as the shape of arguments and return values for various functions. It also contains files that export functions related to email settings and refreshing Google access tokens.

<!--- END SELFDOC --->