package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (i *Infra) setupTopics() error { // create existing gmail pubsub topic
	gmailPubSub, err := pubsub.NewTopic(i.ctx, "gmail-default", &pubsub.TopicArgs{
		Name:    pulumi.String("gmail"),
		Project: pulumi.String(*i.Project.ProjectId),
	}, pulumi.Protect(true))
	if err != nil {
		return err
	}
	i.Topics.DepcrecatedGmail = gmailPubSub

	candidateGmail, err := pubsub.NewTopic(i.ctx, "candidate-gmail", &pubsub.TopicArgs{
		Name:    pulumi.String("candidate-gmail"),
		Project: pulumi.String(*i.Project.ProjectId),
	}, pulumi.Protect(true))
	if err != nil {
		return err
	}
	i.Topics.CandidateGmail = candidateGmail

	recruiterGmail, err := pubsub.NewTopic(i.ctx, "recruiter-gmail", &pubsub.TopicArgs{
		Name:    pulumi.String("recruiter-gmail"),
		Project: pulumi.String(*i.Project.ProjectId),
	}, pulumi.Protect(true))
	if err != nil {
		return err
	}
	i.Topics.RecruiterGmail = recruiterGmail

	return nil
}
