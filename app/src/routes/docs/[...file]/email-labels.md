---
title: Email Labels
description: Learn how SRC uses email labels and folders to manage your inbound recruiting opportunities
---

# {% $frontmatter.title %}

Once you create an account, you'll see a set of new Gmail labels in your inbox. SRC does not store any email information in our database. Instead, we use Gmail labels to add metadata to your existing email inbox. This way you get all the benefits of SRC while keep your data safe and secure.

![SRC Email Labels](/docs/images/gmail-labels.png "Gmail Labels")


## @SRC

{% emailLabel title="@SRC" color="blue" %}{% /emailLabel %}
is the top-level label for all labels managed by SRC. It doesn't serve any functional purpose, but provides a nice folder-like structure for the rest of the labels managed by SRC.

SRC will automatically add this label to any email it takes action on, so it is easily viewable in the Gmail interface.

### @SRC/Jobs

{% emailLabel title="@SRC/Jobs" color="blue" %}{% /emailLabel %}
is a top-level label for all job related emails managed by SRC. It doesn't serve any functional purpose, but provides a nice folder-like structure for the rest of the job labels managed by SRC.

SRC will automatically add this label to any job related email it takes action on, so it is easily viewable in the Gmail interface.

#### @SRC/Jobs/Opportunity

{% emailLabel title="@SRC/Jobs/Opportunity" color="blue" %}{% /emailLabel %}
is attached to any email relating to a job opportunity. SRC defines a job opportunity as an email with a direct contact to someone responsible for that role. 

If SRC miscategorized an email, just remove this label. If you want to help SRC get better, forward any missed or incorrectly labeled emails to [examples@sharedrecruitng.co](mailto:examples@sharedrecruiting.co).

### @SRC/Allow

{% emailLabel title="@SRC/Allow" color="green" %}{% /emailLabel %}
is a top-level label for allowing senders and sender domains to pass through SRC's filter. It doesn't serve any functional purpose, but provides a nice folder-like structure for the rest of the allow list labels used by SRC.

#### @SRC/Allow/Sender

{% emailLabel title="@SRC/Allow/Sender" color="green" %}{% /emailLabel %}
is a user-managed label for allowing specific senders to skip SRC's filter. 

For example, if you are actively talking to a friend about a job opportunity and don't want SRC to manage it, then you can add the
{% emailLabel title="@SRC/Allow/Sender" color="green" %}{% /emailLabel %}
label to any email from your friend and SRC will ignore future emails from this address.

#### @SRC/Allow/Domain

{% emailLabel title="@SRC/Allow/Domain" color="green" %}{% /emailLabel %}
is a user-managed label for allowing specific sender domains (i.e @example.co) to skip SRC's filter. 

A common use-case for 
{% emailLabel title="@SRC/Allow/Domain" color="green" %}{% /emailLabel %}
is allowing all emails from your current organization. For example, say you work at ACME Co. and you know you'll never receive an inbound job opportunity from ACME Co. because you already work there 😉. If you add the 
{% emailLabel title="@SRC/Allow/Domain" color="green" %}{% /emailLabel %}
to _any_ email from someone at you organization, then any future email from your organization will skip SRC's filter.

### @SRC/Block

{% emailLabel title="@SRC/Block" color="red" %}{% /emailLabel %}
is a top-level label for blocking senders and sender domains before they even get to SRC's filter. It doesn't serve any functional purpose, but provides a nice folder-like structure for the rest of the deny list labels used by SRC.

Think of blocking an email as a soft delete. All emails blocked by SRC will be marked as read, archived, and moved to the graveyard ({% emailLabel title="@SRC/Block/🪦" color="red" %}{% /emailLabel %}) for future reference.

#### @SRC/Block/Sender

{% emailLabel title="@SRC/Block/Sender" color="red" %}{% /emailLabel %}
is a user-managed label for blocking specific senders from your inbox. 

Once you label an email with {% emailLabel title="@SRC/Block/Sender" color="red" %}{% /emailLabel %}
, all future emails from that sender will be marked as read, archived, and moved to the blocked graveyard ({% emailLabel title="@SRC/Block/🪦" color="red" %}{% /emailLabel %}) for future reference.

For example, maybe you have a pesky recruiter that is emailing you nonsense trying to get around SRC's filter. To block the recruiter, just add {% emailLabel title="@SRC/Block/Sender" color="red" %}{% /emailLabel %}
 to any email from them and you'll never hear from them again 🤐.

{% emailLabel title="@SRC/Block/Sender" color="red" %}{% /emailLabel %}
 can also be used as a naive unsubscribe from unwanted marketing mailing lists or newsletters.

#### @SRC/Block/Domain

{% emailLabel title="@SRC/Block/Domain" color="red" %}{% /emailLabel %}
 is a user-managed label for blocking specific sender domains (i.e @example.com)

Once you label an email with 
{% emailLabel title="@SRC/Block/Domain" color="red" %}{% /emailLabel %}
 , all future emails from that domain will be marked as read, archived, and moved to the blocked graveyard ({% emailLabel title="@SRC/Block/🪦" color="red" %}{% /emailLabel %}) for future reference.

#### @SRC/Block/🪦

{% emailLabel title="@SRC/Block/🪦" color="red" %}{% /emailLabel %}
is the blocked graveyard. Any email blocked by SRC will be marked as read, archived, and labeled with
{% emailLabel title="@SRC/Block/🪦" color="red" %}{% /emailLabel %}.
The graveyard is here for your peace of mind, you can keep them around or delete them forever. All emails in the graveyard are still searchable in Gmail.

As of right now, emails live in the graveyard until you delete them. In the future, we may add options for automatically pruning emails from your inbox. 
