<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `docs` directory is for the documentation pages related to SRC's recruiting platform. It includes Svelte components for displaying the layout, header, navigation, and error pages. It also contains TypeScript modules and functions for loading and navigating the documentation pages. The directory includes a subdirectory for dynamically loading markdown files and rendering the documentation pages.

### Files
#### +error.svelte
This file contains the Svelte component for displaying an error page. It imports the `PageError` component and renders it with a link to the home page of the documentation section.

#### +layout.svelte
This file contains the layout component for the documentation pages. It imports the Header and Sidenav components and renders them along with a main content area. The Sidenav is hidden on smaller screens and appears as an overlay on larger screens.

#### +page.ts
This file contains a TypeScript module that exports a `load` function. The `load` function redirects the user to the `/docs/welcome` page when the current page is loaded.

#### Header.svelte
This file exports a responsive Svelte component called `Header` that displays a logo, navigation links, and account/login buttons for the documentation pages. It imports the `page` store from the `$app/stores` module and the `Sidenav` component from `./Sidenav.svelte`.

#### Navigation.svelte
This file contains the Navigation component for the documentation pages. It imports the `page` store and the `nav` and `isCurrentPage` functions from the `navigation.ts` file. The component renders a navigation menu with links to different sections of the documentation, highlighting the current page.

#### Sidenav.svelte
This file contains the Sidenav component, which displays a navigation menu on the left side of the screen for the documentation pages. It includes both desktop and mobile versions of the menu, with the mobile menu being hidden by default and revealed when the user clicks on the menu icon. The state of the navigation menu is stored in the isOpen variable, and it uses the fly transition from Svelte and the afterNavigate function from the app/navigation module to close the navigation on mobile devices when the route changes.

#### navigation.ts
This file exports an array of objects representing the navigation links and their titles for the documentation pages. It also exports functions to check if a given page is the current page, get the title of the section a page belongs to, and get the title of a specific page.

### Directories
#### [...file]
The `web/src/routes/docs/[...file]` directory is for documentation pages related to SRC's recruiting platform. It contains files that explain how to contribute to the platform, how to use email settings, how to connect your email to SRC, and how SRC manages inbound job opportunities. It also provides an overview of SRC's open-source policy, security and privacy practices, and the purpose of the platform. The directory includes server-side code for dynamically loading markdown files and a Svelte component for rendering the documentation pages.

<!--- END SELFDOC --->