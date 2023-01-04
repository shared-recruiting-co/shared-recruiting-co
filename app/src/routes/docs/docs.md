# Welcome

Welcome to the SRC documentation!

The Shared Recruiting Co., or SRC (pronounced "source"), is an open source, candidate-centric recruiting platform that promotes two-way opt-in communication between candidates and companies.

For candidates, SRC keeps your inbox distraction free when you aren't looking for a role and supercharges your search once you are. 

For companies, SRC stops you from wasting time sourcing candidates that aren't actively looking for a new role. SRC integrates into your existing recruiting stack and automatically re-engages _interested_ candidates once they are ready for a new role.

**Invite Only** 
SRC is currently invite only. If you know SRC member, have them refer you. If not, join the waitlist. You'll receive an email once you can create an account! 

If you have any questions or just want to chat, feel free to reach out to team@sharedrecruiting.co

Note: At SRC, we use the term candidate loosely. A candidate is both active and prospective candidate. If you are ever looking for a new job, you are a candidate!

# Open Source

## Why is SRC Open Source?

####  Trust 
In the tech world today, trust isn't cheap. Privacy scandal after privacy scandal has left scars on consumers. We believe great companies are built on trust and members of SRC trust us with some of their most sensitive information, their emails and job status. Open sourcing SRC solidifies our commitment to trust, transparency, and privacy and keeps us accountable to our end users.

#### Empower Candidates

Disdain for the current recruiting paradigm is palpable throughout the internet. There are even subreddits dedicated to it. Yet there is no clear way for candidates to take control of their experience.

By open sourcing SRC, we want to empower candidates to create change and have a voice in their experience. For some, that might mean building an integration to make interview scheduling effortless. For others, it could mean starting a discussion on how to disincentivize recruiters from ghosting candidates.

#### Candidates > Code

At the end of the day, SRC's success depends on creating significant, differentiate value for candidates. Code is a means to create that value, but the value will come from candidates themselves.

#### Give Back

Ever wonder how to parse an email body from the Gmail api? Want to know how to forward a message programmatically with the same formatting as your email client? 

While building SRC, we are going to solve many of the same technical problems that others face. Open sourcing our code is another way to give back to the community, share our learnings, and hopefully save some folks a few headaches.

## Can I Contribute?

Yes! There are many ways to contribute to SRC beyond code. Read more in the contributing section.

## Is Any Part of SRC Closed Source?

Yes, the only closed sourced part of SRC is our classification model. At the moment, we believe it's more valuable to candidates to keep the model closed source and make it more difficult for  companies to bypass the model than to make the model open source.  Our hope is to open source this in the future, so will continue to revisit this assumption as we grow.

# Candidate User Guide

## Connect Gmail 

SRC manages your inbound job opportunities for you by connecting to your Gmail account. If you've already created a SRC account, you've already been through this process.

You only need to connect an account once to install SRC. If SRC loses a connection, you'll receive an email to re-connect it. Once re-activated, SRC wile sync your inbox back to the last time your email was active.

You can also pause SRC by clicking the Deactivate button it in your account settings.

**Feature Alert** 
Want to install SRC on multiple emails? Let us know by upvoting this Github issue 

## Labels

Once you create an account, you'll see a set of new Gmail labels in your inbox. SRC does not store any email information in our database. Instead, we use Gmail labels to add metadata to your existing email inbox. This way you get all the benefits of SRC while keep your data safe and secure.

Below is a breakdown of the SRC labels and they are used

Insert screenshot

### @SRC

`@SRC` is the top-level label for all labels managed by SRC. It doesn't serve any functional purpose, but provides a nice folder-like structure for the rest of the labels managed by SRC.

SRC will automatically add this label to any email it takes action on, so it is easily viewable in the Gmail interface.

### @SRC/Jobs

`@SRC/Jobs` is a top-level label for all job related emails managed by SRC. It doesn't serve any functional purpose, but provides a nice folder-like structure for the rest of the job labels managed by SRC.

SRC will automatically add this label to any job related email it takes action on, so it is easily viewable in the Gmail interface.

### @SRC/Jobs/Opportunity

`@SRC/Jobs/Opportunity` is attached to any email relating to a job opportunity. SRC defines a job opportunity as an email with a direct contact to someone responsible for that role. 

If SRC miscategorized an email, just remove this label. If you want to help SRC get better, forward any missed or incorrectly labeled emails to [examples@sharedrecruitng.co](mailto:examples@sharedrecruiting.co)

### @SRC/Allow

`@SRC/Allow`  is a top-level label for allowing senders and sender domains to pass through SRC's filter. It doesn't serve any functional purpose, but provides a nice folder-like structure for the rest of the allow list labels used by SRC.

### @SRC/Allow/Sender

`@SRC/Allow/Sender` is a user-managed label for allowing specific senders to skip SRC's filter. 

For example, if you are actively talking to a friend about a job opportunity and don't want SRC to manage it, then you can add the `@SRC/Allow/Sender` label to any email from your friend and SRC will ignore future emails from this address.

### @SRC/Allow/Domain

`@SRC/Allow/Domain` is a user-managed label for allowing specific sender domains (i.e @example.co) to skip SRC's filter. 

A common use-case for `@SRC/Allow/Domain` is allowing all emails from your current organization. For example, say you work at ACME Co. and you know you'll never receive an inbound job opportunity from ACME Co. because you already work there 😉. If you add the `@SRC/Allow/Domain` to _any_ email from someone at you organization, then any future email from your organization will skip SRC's filter.

### @SRC/Block

`@SRC/Block`  is a top-level label for blocking senders and sender domains before they even get to SRC's filter. It doesn't serve any functional purpose, but provides a nice folder-like structure for the rest of the deny list labels used by SRC.

Think of blocking an email as a soft delete. All emails blocked by SRC will be marked as read, archived, and moved to the graveyard (`@SRC/Block/🪦`) for future reference.

### @SRC/Block/Sender

`@SRC/Block/Sender` is a user-managed label for blocking specific senders from your inbox. 

Once you label an email with `@SRC/Block/Sender` , all future emails from that sender will be marked as read, archived, and moved to the blocked graveyard (`@SRC/Block/🪦`) for future reference.

For example, maybe you have a pesky recruiter that is emailing you nonsense trying to get around SRC's filter. To block the recruiter, just add `@SRC/Block/Sender` to any email from them and you'll never hear from them again 🤐

`@SRC/Block/Sender` can also be used as a naive unsubscribe from unwanted marketing mailing lists or newsletters.

### @SRC/Block/Domain

`@SRC/Block/Domain` is a user-managed label for blocking specific sender domains (i.e @example.com)

Once you label an email with `@SRC/Block/Domain` , all future emails from that domain will be marked as read, archived, and moved to the blocked graveyard (`@SRC/Block/🪦`) for future reference.

### @SRC/Block/🪦

`@SRC/Block/🪦` is the blocked graveyard. Any email blocked by SRC will be marked as read, archived, and labeled with `@SRC/Block/🪦` . The graveyard is here for your peace of mind, you can keep them around or delete them forever. All emails in the graveyard are still searchable in Gmail.

As of right now, emails live in the graveyard until you delete them. In the future, we may add options for automatically pruning emails from your inbox. 

## Email Settings

SRC has a number of email settings to tune SRC to your preferences.

### Hide Recruiting Emails from Inbox

Keep your inbox distraction free when you aren't actively looking for a new role.

When this setting is enabled, all recruiting emails will be automatically marked as read, archived, and moved under the `@SRC/Job/Opportunity` folder.

### Block Automated Email Sequences  (Coming Soon)

Block automated recruiting sequences by automatically replying to recruiters with a standard message.

When this setting is enabled, all new recruiting emails will be receive an automated reply to prevent automated follow ups.

**Coming Soon**
Want to try out this feature? Let us know by 👍 this issue on Github

### Auto-Contribute Recruiting Emails

At SRC, we want to build the best recruiting experience for candidates. We want to automate away all the parts of the recruiting process that isn't essential to the qualification process.

To do this, we are collecting examples of recruiting emails to train our models on. When this settings is enabled, all recruiting emails you receive will be added to the recruiting email dataset. These emails will exclusively used for improving SRC and will not be shared with anyone.

You can always contribute emails manually by forwarding them to examples@sharedrecruiting.co. 

# Security and Privacy

## Open Source

We believe [Linus's law](https://en.wikipedia.org/wiki/Linus%27s_law) of "given enough eyeballs, all bugs are shallow" applies to security issues. The premier example of how open source projects can be more secure than proprietary code bases is Bitcoin. In [his 2015 talk](https://www.youtube.com/watch?v=810aKcfM__Q) Andreas M. Antonopoulos describes how closed source banking systems have the software equivalent of weak immune systems, because huge security holes can be obfuscated for long periods of time, and when eventually exploited can have enormous detrimental effects. On the flip side of this is an open source protocol like Bitcoin, where any security holes are there for all to see. Exploits are found early and often, and then patched. Remember that successful software companies can take more than a decade to build. Over a long time period, open source systems will tend towards a more secure state over secretive, proprietary systems.

## Privacy-First

At SRC, members trust us with some of their most sensitive information, their emails and their job status. Sadly, we've seen firsthand what happens when companies abuse this.

We've seen employers use LinkedIn to find out which employees are looking for new roles. We've seen backlash against recruiting companies for making candidate job profiles public without their permission.

SRC does not and will never store your emails within our database. We use native email functionality, like labels and folders, to securely manage your emails from within your inbox.  
And with SRC's two-way opt-in communication, your job status is confidential. Companies cannot see if you are looking for a new role until you choose to start the interview process.

By open sourcing SRC, we are making more then a promise of privacy. All the ways we use and protect your data will is transparent and publicly visible.

## CASA Compliant

The Cloud Application Security Assessment (CASA) is built upon the industry-recognized standards of the [OWASP's Application Security Verification Standard (ASVS)](https://owasp.org/www-project-application-security-verification-standard/) to provide a consistent set of requirements to harden security for any application. CASA provides a uniform way to perform trusted assurance assessments of these requirements when such assessments are required for applications with potential access to sensitive data.

SRC is Tier-2 CASA compliant. SRC's CASA compliance has been verified by an independent lab partner as part of the Google OAuth approval process. 

## Privacy Policy & Terms of Service

SRC's privacy policy and terms of service are available on the SRC website. Like the rest of SRC, the privacy policy and terms of service are open source too.

We want to make the privacy policy and terms of service as interpretable and user-friendly as possible. If you find the language confusing or ambiguous, please let us know!

# Contributing

## Community

You don't have to know how to code to contribute to SRC. The best way to contribute to SRC is by participating in the community!

- [Share and upvote features](https://github.com/shared-recruiting-co/shared-recruiting-co/discussions/categories/ideas)
- [Post and answer questions](https://github.com/shared-recruiting-co/shared-recruiting-co/discussions/categories/q-a)
- [Poll the community](https://github.com/shared-recruiting-co/shared-recruiting-co/discussions/categories/polls)

## Code

Found a typo on the website? Want to add a cool animation? Know how to improve the performance of one of SRC's cloud functions? 

Fork the repository, make a PR, and contribute! You can read more in our contribution guidelines.

## Architecture

To help you get familiar with the codebase and how SRC works, here is the high-level architecture

![SRC Architecture Diagram](https://github.com/shared-recruiting-co/shared-recruiting-co/raw/main/architecture.png "Architecture")

## Emails

One of the simplest ways to contribute to SRC is by sharing recruiting emails. At SRC, we want to build the best recruiting experience for candidates. We want to automate away all the parts of the recruiting process that isn't essential to the qualification process.

To do this, we are collecting recruiting emails to train our models on. If SRC misses a recruiting email or labels an email incorrectly, help us improve by forwarding the email to [examples@sharedrecruiting.co](mailto:examples@sharedrecruiting.co).
