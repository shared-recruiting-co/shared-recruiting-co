package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (i *Infra) setupTopics() error {
	gmailPubSub, err := pubsub.NewTopic(i.ctx, "gmail-default", &pubsub.TopicArgs{
		Name:    pulumi.String("gmail"),
		Project: pulumi.String(*i.Project.ProjectId),
	}, pulumi.Protect(true))
	if err != nil {
		return err
	}
	i.Topics.Gmail = gmailPubSub

	candidateGmailMessages, err := pubsub.NewTopic(i.ctx, "candidate-gmail-messages", &pubsub.TopicArgs{
		Name:    pulumi.String("candidate-gmail-messages"),
		Project: pulumi.String(*i.Project.ProjectId),
		// TODO: Enforce schema validation
	}, pulumi.Protect(true))
	if err != nil {
		return err
	}
	i.Topics.CandidateGmailMessages = candidateGmailMessages

	return nil
}