package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	serviceAccount "github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	MaxEventArcTriggerTimeout = 540
	MaxHTTPTriggerTimeout     = 3600
	DefaultRegion             = "us-west1"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
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
		_, err = pubsub.NewTopic(ctx, "gmail-default", &pubsub.TopicArgs{
			Name:    pulumi.String("gmail"),
			Project: pulumi.String(*project.ProjectId),
		}, pulumi.Protect(true))
		if err != nil {
			return err
		}

		_, err = storage.NewBucketObject(ctx, "cf-upload-email-push-notifications", &storage.BucketObjectArgs{
			Name:   pulumi.String("email-push-notifications/function.zip"),
			Bucket: gcfBucket.Name,
			Source: pulumi.NewFileArchive("../cloudfunctions/email_push_notifications"),
		}, pulumi.DependsOn([]pulumi.Resource{
			gcfBucket,
		}))
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
		sa, err := serviceAccount.NewAccount(ctx, "sa-cf-email-push-notification", &serviceAccount.AccountArgs{
			Project:     pulumi.String(*project.ProjectId),
			AccountId:   pulumi.String("sa-cf-email-push-notification"),
			DisplayName: pulumi.Sprintf("Service account for the %s Cloud Function", "email-push-notifications"),
		})
		if err != nil {
			return err
		}
		_, err = projects.NewIAMMember(ctx, "sa-cf-email-push-notifications-invoker", &projects.IAMMemberArgs{
			Project: pulumi.String(*project.ProjectId),
			Role:    pulumi.String("roles/run.invoker"),
			Member: sa.Email.ApplyT(func(email string) (string, error) {
				return fmt.Sprintf("serviceAccount:%v", email), nil
			}).(pulumi.StringOutput),
		}, pulumi.DependsOn([]pulumi.Resource{
			sa,
		}))
		if err != nil {
			return err
		}

		// allow the service account to create ID tokens to authenticate with other services
		_, err = projects.NewIAMMember(ctx, "sa-cf-email-push-notifications-open-id-token-creator", &projects.IAMMemberArgs{
			Project: pulumi.String(*project.ProjectId),
			Role:    pulumi.String("roles/iam.serviceAccountOpenIdTokenCreator"),
			Member: sa.Email.ApplyT(func(email string) (string, error) {
				return fmt.Sprintf("serviceAccount:%v", email), nil
			}).(pulumi.StringOutput),
		}, pulumi.DependsOn([]pulumi.Resource{
			sa,
		}))
		if err != nil {
			return err
		}
		// _, err = cloudfunctionsv2.NewFunction(ctx, "testing-pulumi", &cloudfunctionsv2.FunctionArgs{
		// // use the same location as the bucket
		// Name:        pulumi.Sprintf("test-%s", DefaultRegion),
		// Location:    pulumi.String(DefaultRegion),
		// Description: pulumi.String("a new function"),
		// BuildConfig: &cloudfunctionsv2.FunctionBuildConfigArgs{
		// Runtime:    pulumi.String("go119"),
		// EntryPoint: pulumi.String("Handler"),
		// EnvironmentVariables: pulumi.StringMap{
		// "BUILD_CONFIG_TEST": pulumi.String("build_test"),
		// },
		// Source: &cloudfunctionsv2.FunctionBuildConfigSourceArgs{
		// StorageSource: &cloudfunctionsv2.FunctionBuildConfigSourceStorageSourceArgs{
		// Bucket: gcfBucket.Name,
		// Object: cf.Name,
		// },
		// },
		// },
		// ServiceConfig: &cloudfunctionsv2.FunctionServiceConfigArgs{
		// MaxInstanceCount: pulumi.Int(1),
		// MinInstanceCount: pulumi.Int(1),
		// AvailableMemory:  pulumi.String("256M"),
		// TimeoutSeconds:   pulumi.Int(MaxEventArcTriggerTimeout),
		// EnvironmentVariables: pulumi.StringMap{
		// "SERVICE_CONFIG_TEST": pulumi.String("config_test"),
		// },
		// IngressSettings:            pulumi.String("ALLOW_INTERNAL_ONLY"),
		// AllTrafficOnLatestRevision: pulumi.Bool(true),
		// ServiceAccountEmail:        sa.Email,
		// },
		// EventTrigger: &cloudfunctionsv2.FunctionEventTriggerArgs{
		// TriggerRegion: pulumi.String(DefaultRegion),
		// PubsubTopic:   gmailPubSub.ID(),
		// EventType:     pulumi.String("google.cloud.pubsub.topic.v1.messagePublished"),
		// // Disable retry
		// RetryPolicy:         pulumi.String("RETRY_POLICY_DO_NOT_RETRY"),
		// ServiceAccountEmail: sa.Email,
		// },
		// }, pulumi.DependsOn([]pulumi.Resource{
		// gmailPubSub,
		// cf,
		// sa,
		// }))
		// if err != nil {
		// return err
		// }
		return nil
	})
}
