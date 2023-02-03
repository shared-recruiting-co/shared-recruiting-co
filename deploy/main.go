package main

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudfunctionsv2"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	serviceAccount "github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
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

type Infra struct {
	ctx       *pulumi.Context
	config    *config.Config
	ProjectID string
	Project   *organizations.LookupProjectResult
	GCFBucket *storage.Bucket
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
			ProjectID: *project.ProjectId,
			Project:   project,
			GCFBucket: gcfBucket,
		}

		// create existing gmail pubsub topic
		gmailPubSub, err = pubsub.NewTopic(ctx, "gmail-default", &pubsub.TopicArgs{
			Name:    pulumi.String("gmail"),
			Project: pulumi.String(*project.ProjectId),
		}, pulumi.Protect(true))
		if err != nil {
			return err
		}

		infra.emailPushNotificationCF(gmailPubSub)

		return nil
	})
}

func (i *Infra) uploadCloudFunction(name string) (*storage.BucketObject, error) {
	snakeName := strings.ReplaceAll(name, "-", "_")
	// for now hardcode path
	dir := "../cloudfunctions"

	return storage.NewBucketObject(i.ctx, fmt.Sprintf("cf-upload-%s", name), &storage.BucketObjectArgs{
		Name:   pulumi.Sprintf("%s/function.zip", name),
		Bucket: i.GCFBucket.Name,
		Source: pulumi.NewFileArchive(fmt.Sprintf("%s/%s", dir, snakeName)),
	}, pulumi.DependsOn([]pulumi.Resource{
		i.GCFBucket,
	}))
}

func (i *Infra) createCloudFunctionServiceAccount(name string) (*serviceAccount.Account, error) {
	account := fmt.Sprintf("sa-cf-%s", name)

	sa, err := serviceAccount.NewAccount(i.ctx, account, &serviceAccount.AccountArgs{
		Project:     pulumi.String(*i.Project.ProjectId),
		AccountId:   pulumi.String(account),
		DisplayName: pulumi.Sprintf("Service account for the %s Cloud Function", name),
	})
	if err != nil {
		return nil, err
	}
	_, err = projects.NewIAMMember(i.ctx, fmt.Sprintf("%s-invoker", account), &projects.IAMMemberArgs{
		Project: pulumi.String(*i.Project.ProjectId),
		Role:    pulumi.String("roles/run.invoker"),
		Member: sa.Email.ApplyT(func(email string) (string, error) {
			return fmt.Sprintf("serviceAccount:%v", email), nil
		}).(pulumi.StringOutput),
	}, pulumi.DependsOn([]pulumi.Resource{
		sa,
	}))
	if err != nil {
		return nil, err
	}

	// allow the service account to create ID tokens to authenticate with other services
	_, err = projects.NewIAMMember(i.ctx, fmt.Sprintf("%s-open-id-token-creator", account), &projects.IAMMemberArgs{
		Project: pulumi.String(*i.Project.ProjectId),
		Role:    pulumi.String("roles/iam.serviceAccountOpenIdTokenCreator"),
		Member: sa.Email.ApplyT(func(email string) (string, error) {
			return fmt.Sprintf("serviceAccount:%v", email), nil
		}).(pulumi.StringOutput),
	}, pulumi.DependsOn([]pulumi.Resource{
		sa,
	}))
	if err != nil {
		return nil, err
	}
	return sa, nil
}

func (i *Infra) emailPushNotificationCF(topic *pubsub.Topic) error {
	name := "email-push-notifications"
	sa, err := i.createCloudFunctionServiceAccount(name)
	if err != nil {
		return err
	}
	obj, err := i.uploadCloudFunction(name)
	if err != nil {
		return err
	}

	_, err = cloudfunctionsv2.NewFunction(i.ctx, fmt.Sprintf("%s-%s", name, DefaultRegion), &cloudfunctionsv2.FunctionArgs{
		// use the same location as the bucket
		Name:        pulumi.Sprintf("%s-%s", name, DefaultRegion),
		Location:    pulumi.String(DefaultRegion),
		Description: pulumi.String("Handle user email push notifications"),
		BuildConfig: &cloudfunctionsv2.FunctionBuildConfigArgs{
			Runtime:    pulumi.String("go119"),
			EntryPoint: pulumi.String("EmailPushNotificationsHandler"),
			Source: &cloudfunctionsv2.FunctionBuildConfigSourceArgs{
				StorageSource: &cloudfunctionsv2.FunctionBuildConfigSourceStorageSourceArgs{
					Bucket: i.GCFBucket.Name,
					Object: obj.Name,
				},
			},
		},
		ServiceConfig: &cloudfunctionsv2.FunctionServiceConfigArgs{
			AvailableMemory:  pulumi.String("256M"),
			MinInstanceCount: pulumi.Int(1),
			MaxInstanceCount: pulumi.Int(25),
			TimeoutSeconds:   pulumi.Int(MaxEventArcTriggerTimeout),
			EnvironmentVariables: pulumi.StringMap{
				"SUPABASE_API_URL":           pulumi.String(i.config.Require("SUPABASE_API_URL")),
				"SUPABASE_API_KEY":           i.config.RequireSecret("SUPABASE_API_KEY"),
				"GOOGLE_OAUTH2_CREDENTIALS":  i.config.RequireSecret("GOOGLE_OAUTH2_CREDENTIALS"),
				"ML_SERVICE_URL":             i.config.RequireSecret("ML_SERVICE_URL"),
				"SENTRY_DSN":                 i.config.RequireSecret("SENTRY_DSN"),
				"EXAMPLES_GMAIL_OAUTH_TOKEN": i.config.RequireSecret("EXAMPLES_GMAIL_OAUTH_TOKEN"),
				// TOOD: Parameterize this
				"TRIGGER_FULL_SYNC_URL": pulumi.String("https://full-email-sync-bjrwwezbha-uw.a.run.app"),
			},
			IngressSettings:            pulumi.String("ALLOW_INTERNAL_ONLY"),
			AllTrafficOnLatestRevision: pulumi.Bool(true),
			ServiceAccountEmail:        sa.Email,
		},
		EventTrigger: &cloudfunctionsv2.FunctionEventTriggerArgs{
			TriggerRegion: pulumi.String(DefaultRegion),
			PubsubTopic:   topic.ID(),
			EventType:     pulumi.String("google.cloud.pubsub.topic.v1.messagePublished"),
			// Disable retry
			RetryPolicy:         pulumi.String("RETRY_POLICY_DO_NOT_RETRY"),
			ServiceAccountEmail: sa.Email,
		},
	}, pulumi.DependsOn([]pulumi.Resource{
		topic,
		obj,
		sa,
	}))

	if err != nil {
		return err
	}

	return nil
}
