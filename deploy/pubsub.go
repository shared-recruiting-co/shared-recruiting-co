package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (i *Infra) setupTopics() error {
	candidateGmailSubscription, err := pubsub.NewTopic(i.ctx, "candidate-gmail-subscription", &pubsub.TopicArgs{
		Name:    pulumi.String("candidate-gmail-subscription"),
		Project: pulumi.String(*i.Project.ProjectId),
		// TODO: Enforce schema validation
	}, pulumi.Protect(true))
	if err != nil {
		return err
	}
	i.Topics.CandidateGmailSubscription = candidateGmailSubscription

	recruiterGmailSubscription, err := pubsub.NewTopic(i.ctx, "recruiter-gmail-subscription", &pubsub.TopicArgs{
		Name:    pulumi.String("recruiter-gmail-subscription"),
		Project: pulumi.String(*i.Project.ProjectId),
		// TODO: Enforce schema validation
	}, pulumi.Protect(true))
	if err != nil {
		return err
	}
	i.Topics.RecruiterGmailSubscription = recruiterGmailSubscription

	candidateGmailMessages, err := pubsub.NewTopic(i.ctx, "candidate-gmail-messages", &pubsub.TopicArgs{
		Name:    pulumi.String("candidate-gmail-messages"),
		Project: pulumi.String(*i.Project.ProjectId),
		// TODO: Enforce schema validation
	}, pulumi.Protect(true))
	if err != nil {
		return err
	}
	i.Topics.CandidateGmailMessages = candidateGmailMessages

	recruiterGmailMessages, err := pubsub.NewTopic(i.ctx, "recruiter-gmail-messages", &pubsub.TopicArgs{
		Name:    pulumi.String("recruiter-gmail-messages"),
		Project: pulumi.String(*i.Project.ProjectId),
		// TODO: Enforce schema validation
	}, pulumi.Protect(true))
	if err != nil {
		return err
	}
	i.Topics.RecruiterGmailMessages = recruiterGmailMessages

	candidateGmailLabelChanges, err := pubsub.NewTopic(i.ctx, "candidate-gmail-label-changes", &pubsub.TopicArgs{
		Name:    pulumi.String("candidate-gmail-label-changes"),
		Project: pulumi.String(*i.Project.ProjectId),
		// TODO: Enforce schema validation
	}, pulumi.Protect(true))
	if err != nil {
		return err
	}
	i.Topics.CandidateGmailLabelChanges = candidateGmailLabelChanges

	return nil
}
