<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `components` directory is for exporting Svelte components that are used in the SRC platform. These components include a modal, a Google account connection button, a hint icon with a tooltip, a page error message, a table footer with pagination controls, and a customizable toggle switch. It also includes directories for rendering markdown and marketing-related content.

### Files
#### AlertModal.svelte
This file (`AlertModal.svelte`) contains a Svelte component that displays a modal with a title, description, and call-to-action button. The modal fades in and out using a transition effect. The component also includes a function to close the modal and a callback function to execute when the call-to-action button is clicked. The modal also includes an exclamation-triangle icon and has an animation that slides up and down.

#### ConnectGoogleAccountButton.svelte
This file exports a Svelte component called `ConnectGoogleAccountButton` that allows users to connect their Google account to the SRC platform. It uses the Google OAuth2 API to authenticate the user and make a POST request to the `/api/account/gmail/connect` endpoint. The component handles errors and loading states and has props for email, onConnect, and disabled.

#### Hint.svelte
This file exports a Svelte component that displays a hint icon and a tooltip when the user hovers over the icon. The tooltip content is passed in as a prop to the component.

#### PageError.svelte
This file defines the PageError component, which displays an error message and a link to return to the home page or contact the team. The error message is either "Page not found" or "Something went wrong" depending on the error status, and a message with more details is displayed if available.

#### TableFooter.svelte
This file exports a Svelte component called `TableFooter` that displays pagination controls and result count information at the bottom of a table. It includes "Previous" and "Next" buttons, as well as links to individual pages. The appearance of the links changes depending on whether they represent the current page or not. The component receives a `pagination` object as a prop and imports a `Pagination` type and a `getUpdatedPageUrl` function from `$lib/pagination`.

#### Toggle.svelte
This file contains the `Toggle` component, which is a customizable toggle switch with a label, a boolean value for whether it is checked or not, and an optional callback function to be called when the toggle is switched. It also has an optional disabled state and includes accessibility features such as aria-checked and role attributes.

### Directories
#### markdoc
The `markdoc` directory is for exporting Svelte components and configuration objects used in rendering markdown as HTML using the `@markdoc/markdoc` library. It includes components for rendering styled boxes, email labels, and tags, as well as a configuration file for generating a configuration object based on an abstract syntax tree (AST).

#### marketing
The `marketing` directory is for Svelte components that render marketing-related content, such as FAQs, features, and a step-by-step guide on how the Shared Recruiting Co. platform works.

<!--- END SELFDOC --->