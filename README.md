# Shared Recruiting Co. (SRC)

[![CI](https://github.com/shared-recruiting-co/shared-recruiting-co/actions/workflows/ci.yml/badge.svg)](https://github.com/shared-recruiting-co/shared-recruiting-co/actions/workflows/ci.yml) [![CodeQL](https://github.com/shared-recruiting-co/shared-recruiting-co/actions/workflows/codeql.yml/badge.svg)](https://github.com/shared-recruiting-co/shared-recruiting-co/actions/workflows/codeql.yml) [![Database Migrations](https://github.com/shared-recruiting-co/shared-recruiting-co/actions/workflows/migrations.yml/badge.svg)](https://github.com/shared-recruiting-co/shared-recruiting-co/actions/workflows/migrations.yml)

Welcome to the SRC monorepo 👋

The Shared Recruiting Company, SRC (pronounced "source"), is an open source, candidate-centric recruiting platform. SRC promotes two-way, opt-in communication between candidates and companies. 🤝

For candidates, SRC is a recruiting AI assistant that lives in your inbox. No more recruiting spam emails vying for your attention. SRC manages your inbound job opportunities when you aren't looking for a new role and supercharges your job search once you are.

For companies, SRC stops you from wasting time sourcing candidates that aren't actively looking for a new role. SRC integrates into your existing recruiting stack and automatically re-engages interested candidates once they are ready for a new role.

## 😎 Become a Member

Right now, SRC is invite only. If you are interesting in joining, sign up at [sharedrecruiting.co](https://sharedrecruiting.co/).

## 📖 Documentation

For user-facing app documentation, checkout the [SRC Docs](https://sharedrecruiting.co/docs/welcome).

For code-related documentation, all the documentation lives in this repository. Read the directory-level README, read the code comments, or read the code itself 🤓

## 🕍 Project Layout

### 📱 `/web`

The SRC web app ([sharedrecruiting.co](https://sharedrecruiting.co)) built via Sveltekit + Tailwind + Supabase and deployed with Vercel

#### Development

First setup your local environment variables. From the `/web` directory, run 

```bash
cp .env.example .env.local
# Now edit .env.local and replace the values to match your local setup
# Use whichever editor you prefer (I use vim)
vi .env.local
```



To start the web app, run
```bash
npm run dev -- --open 
```

To run a local instance of Supabase, run
```bash
supabase start
```

If you want to log into the app locally, add your Google OAuth client ID/secret to the bottom of `web/supabase/config.toml`:

```toml
[auth.external.google]
enabled = true
client_id = "xxx"
secret = "xxx"
```

Changes to `web/supabase/config.toml` are ignored by Git, so you don't have to worry about accidentally committing your client secret.

#### Testing

Browser and Typescript and test are written via [Playwright](https://playwright.dev/) and [Vitest](https://vitest.dev/). You can run them with,

```bash
npm run test
```

Database tests are powered by [pgtap](https://pgtap.org) and allow us to test complex database logic, like RLS policies and triggers. The tests live in `web/supabase/tests/database`. You can run them with

```bash
supabase test db
```

Note: The `web/` is under active development. Test coverage is intentionally low until the app stabilizes. 

### 🌩️ `/cloudfunctions`

The SRC Google Cloud Functions. The cloud functions are responsible for managing and reacting to user emails. To minimize unnecessary dependencies, each cloud function is an independent, deployable  `go` module. 

### 🎮 `/libs`

Shared `go` libraries

### 📑 `/scripts`

Scripts for common manual tasks

## 👩‍💻 Contributing

SRC is open source to empower candidates and companies to contribute and collaborate on building an ideal and efficient recruiting experience for all.

Have a feature idea? Create an [issue](https://github.com/shared-recruiting-co/shared-recruiting-co/issues). Want to fix a bug? Create a [pull request](https://github.com/shared-recruiting-co/shared-recruiting-co/pulls). Have a question? Start a [discussion](https://github.com/shared-recruiting-co/shared-recruiting-co/discussions).

### Contribute Recruiting Emails

We want to build the best candidates experience possible. To do so, SRC needs examples of all types of inbound recruiting emails. If you have inbound recruiting emails you want to contribute to our dataset please forward them to [examples@sharedrecruiting.co](mailto:examples@sharedrecruiting.co) 

## 🖼️ Architecture

![SRC Architecture Diagram](/web/static/docs/images/architecture.png "Architecture")


<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `root` directory is for managing the high-level documentation and configuration files of the Shared Recruiting Co. (SRC) platform. It includes files for managing contributions, code of conduct, security, and Git tracking. Additionally, it includes directories for managing Cloud Functions, deployment, shared libraries, scripts, and the web application.

### Files
#### .gitignore
This file contains a list of files and directories that should be ignored by Git when tracking changes in the codebase. The listed files and directories include environment files and macOS-specific files.

#### CODE_OF_CONDUCT.md
This file contains the SRC Code of Conduct, which outlines the expected behavior of community members and contributors, and the consequences for violating those expectations. It includes a pledge to make participation in the community a harassment-free experience for everyone and examples of acceptable and unacceptable behavior. Community leaders are responsible for enforcing these standards.

#### CONTRIBUTING.md
This file is a comprehensive guide for contributing to the SRC project. It includes instructions on forking the repository, making changes locally, creating and solving issues, creating a pull request, and enabling maintainer edits. It also congratulates contributors on their merged PR and mentions that their contributions will be deployed in the next release.

#### SECURITY.md
This file contains the security policy for SRC. It outlines the currently supported version and provides instructions for reporting a vulnerability.

#### fluid-config.yaml
This file contains the configuration settings for SRC's security analysis tool and language used in the platform. It includes options for specifying the namespace, output format, working directory, target files, and findings to analyze. The default language is English, but it can be changed to Spanish by modifying the "language" field in this YAML file. Additionally, the file provides a link to the complete list of findings that can be analyzed.

### Directories
#### .github
The `.github` directory is for managing various aspects of the SRC codebase on GitHub. It includes a `FUNDING.yml` file for listing supported funding platforms, a `dependabot.yml` file for configuring Dependabot, a pull request template, and directories for issue templates and workflows. These files and directories help streamline contributions, automate tasks, and ensure code quality.

#### cloudfunctions
The `cloudfunctions` directory is for various Cloud Functions that handle different tasks related to the Shared Recruiting Co. platform. These tasks include syncing candidate and recruiter emails, handling changes to Gmail labels, processing and labeling recruiting emails, and handling Gmail push notifications. Each subdirectory contains implementation files, test functions, configuration files, and module dependencies.

#### deploy
The `deploy` directory is for deploying the Shared Recruiting Co. platform to Google Cloud Platform using Pulumi and Go code. It includes configuration files for Pulumi and GolangCI, Go code for creating and deploying various Google Cloud Functions, and files for managing dependencies and testing.

#### libs
The `libs` directory is for managing the configuration and source code of the Shared Recruiting Co. platform. It includes packages for interacting with the SRC database, Gmail API, and machine learning functionality related to recruiting emails and job postings. Additionally, it includes a directory for implementing the pub/sub system in the SRC platform.

#### scripts
The `scripts` directory is for storing scripts related to the development and maintenance of the SRC codebase. This includes a script called `update-go-mods.sh` which updates the dependencies of several Go packages in the `cloudfunctions` directory and runs tests.

#### web
The `web` directory is for managing the web application of the SRC (Shared Recruiting Co.) platform. It includes configuration files for ESLint, Prettier, PostCSS, TypeScript, and Vite. It also contains the `src` directory for managing the source code, the `static` directory for storing static assets, the `supabase` directory for managing the Supabase project configuration, and the `tests` directory for storing Playwright tests. Additionally, it includes a `README.md` file with instructions for creating, developing, and building a Svelte project.

<!--- END SELFDOC --->