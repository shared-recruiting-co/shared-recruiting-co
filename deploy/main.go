package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

var (
	MaxEventArcTriggerTimeout = 540
	MaxHTTPTriggerTimeout     = 3600
	DefaultRegion             = "us-west1"
)

// TODOs
// 3. How do we grant invoker permission to service account to execute other cloudfunction's cloud run service
// 5. Cloud Scheduler
// Read about https://www.pulumi.com/docs/guides/testing/

type Topics struct {
	DepcrecatedGmail *pubsub.Topic
	CandidateGmail   *pubsub.Topic
	RecruiterGmail   *pubsub.Topic
}

type Infra struct {
	ctx       *pulumi.Context
	config    *config.Config
	Project   *organizations.LookupProjectResult
	GCFBucket *storage.Bucket
	Topics    Topics
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		// Get current project ID
		project, err := organizations.LookupProject(ctx, &organizations.LookupProjectArgs{})
		if err != nil {
			return err
		}
		// Bucket names must be globally unique
		gcfBucket, err := storage.NewBucket(ctx, "gcf-uploads", &storage.BucketArgs{
			Name:                     pulumi.String(fmt.Sprintf("gcf-uploads-%s-%s", project.Number, DefaultRegion)),
			Location:                 pulumi.String(DefaultRegion),
			UniformBucketLevelAccess: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		// Grant publish permission to the gmail topic
		// https://developers.google.com/gmail/api/guides/push
		_, err = projects.NewIAMMember(ctx, "gmail-pubsub-publishing", &projects.IAMMemberArgs{
			Project: pulumi.String(*project.ProjectId),
			Role:    pulumi.String("roles/pubsub.publisher"),
			Member:  pulumi.String("serviceAccount:gmail-api-push@system.gserviceaccount.com"),
		})
		if err != nil {
			return err
		}

		infra := &Infra{
			ctx:       ctx,
			config:    cfg,
			Project:   project,
			GCFBucket: gcfBucket,
		}

		err = infra.setupTopics()
		if err != nil {
			return err
		}

		err = infra.createCloudFunctions()
		if err != nil {
			return err
		}

		return nil
	})
}
