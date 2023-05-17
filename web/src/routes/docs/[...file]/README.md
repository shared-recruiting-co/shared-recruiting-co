<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `web/src/routes/docs/[...file]` directory is for documentation pages related to SRC's recruiting platform. It contains files that explain how to contribute to the platform, how to use email settings, how to connect your email to SRC, and how SRC manages inbound job opportunities. It also provides an overview of SRC's open-source policy, security and privacy practices, and the purpose of the platform. The directory includes server-side code for dynamically loading markdown files and a Svelte component for rendering the documentation pages.

### Files
#### +page.server.ts
This file is the server-side code for dynamically loading markdown files for the documentation pages. It redirects to the welcome page if the requested file does not exist.

#### +page.svelte
This file contains the Svelte component for rendering a documentation page. It imports and uses the `Markdoc` component to render the markdown content of the page. It also sets the page title and section title based on the URL.

#### contributing.md
This file explains how to contribute to SRC, including participating in the community, contributing code, and improving the architecture. It also provides a high-level architecture diagram and encourages users to share recruiting emails to help train SRC's models.

#### email-labels.md
This file explains how SRC uses Gmail labels and folders to manage inbound recruiting opportunities. It lists the different labels used by SRC, including @SRC, @SRC/Jobs, and @SRC/Jobs/Opportunity, and explains their purposes. It also provides instructions on how to remove or forward miscategorized emails, as well as how to use the email blocking feature with user-managed labels. Additionally, it explains the purpose of the "@SRC/Block/🪦" email label and the possibility of automatic email pruning in the future.

#### email-settings.md
This file explains how to use email settings to customize SRC's email behavior to your preferences. It covers three settings: hiding recruiting emails from your inbox, blocking automated email sequences (coming soon), and auto-contributing recruiting emails to SRC's dataset to improve the platform.

#### email-setup.md
This file explains how to connect your email to SRC. SRC manages your inbound job opportunities by syncing your inbox. You only need to connect your account once and can pause SRC at any time.

#### open-source.md
This file provides an overview of SRC's open-source policy, including the benefits of trust, empowering candidates, and creating value for candidates. It also explains how to contribute to the project beyond code and mentions the only closed-sourced part of SRC.

#### security-privacy.md
This file provides an overview of SRC's security and privacy practices, including its commitment to privacy-first practices and the benefits of open source projects for security. It also highlights the importance of protecting sensitive information and provides examples of how it can be abused. Additionally, it mentions SRC's Tier-2 CASA compliance and provides links to its privacy policy and terms of service, which are open source and user-friendly.

#### welcome.md
This file explains the purpose of the SRC platform and its benefits for both candidates and companies. It also provides information on how to join the platform, and defines the term "candidate" as used by SRC.

<!--- END SELFDOC --->