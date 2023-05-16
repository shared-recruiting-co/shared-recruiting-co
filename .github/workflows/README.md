<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `workflows` directory is for storing GitHub Actions workflows that automate various tasks such as auto-assigning pull requests, running CI tests, analyzing code with CodeQL, and deploying migrations to production. These workflows are triggered on specific events and can be customized to fit the needs of the SRC codebase.

### Files
#### auto-assign.yml
This file contains a GitHub Actions workflow that automatically assigns pull requests to reviewers and/or assignees based on specified criteria. It uses the `wow-actions/auto-assign` action and allows for adding reviewers by team or individual, setting the number of reviewers, and skipping certain keywords and labels.

#### ci.yml
This file contains the CI workflow for the SRC codebase, which runs tests, lints the codebase using golangci-lint, and builds the Go code in various directories. It also includes a job for running Supabase tests and verifying generated types.

#### codeql.yml
This file configures a GitHub Actions workflow that runs CodeQL analysis on the repository's codebase. It is triggered on push and pull request events on the main branch, as well as on a weekly schedule. The workflow is set up to analyze Go and TypeScript code, but this can be customized. It also includes an autobuild step for compiled languages and a CodeQL analysis step, with instructions for manual building if the autobuild step fails.

#### migrations.yml
This file contains a workflow that deploys migrations to production. It runs on the main branch and allows manual triggers. It uses Supabase CLI to push database changes to production.

<!--- END SELFDOC --->