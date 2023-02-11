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

### 📱 `/app`

The SRC web app ([sharedrecruiting.co](https://sharedrecruiting.co)) built via Sveltekit + Tailwind + Supabase and deployed with Vercel

#### Development

To start the web app, run
```bash
npm run dev -- --open 
```

To run a local instance of Supabase, run
```bash
supabase start
```

If you want to log into the app locally, add your Google OAuth client ID/secret to the bottom of `app/supabase/config.toml`:

```toml
[auth.external.google]
enabled = true
client_id = "xxx"
secret = "xxx"
```

Changes to `app/supabase/config.toml` are ignored by Git, so you don't have to worry about accidentally committing your client secret.

#### Testing

Tests are written via [Playwright](https://playwright.dev/) and [Vitest](https://vitest.dev/). You can run them via,

```bash
npm run test
```

Note: The `app/` is under active development. Test coverage is intentionally low until the app stabilizes. 

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

![SRC Architecture Diagram](/app/static/docs/images/architecture.png "Architecture")
