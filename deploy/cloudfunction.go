package main

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudfunctionsv2"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudrun"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudrunv2"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	serviceAccount "github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CloudFunction struct {
	Name           string
	ServiceAccount *serviceAccount.Account
	Function       *cloudfunctionsv2.Function
	Service        *cloudrun.LookupServiceResult
}

func (i *Infra) createCloudFunctions(gmailPubSub *pubsub.Topic) error {
	syncCF, err := i.fullEmailSyncCF()
	if err != nil {
		return err
	}
	i.ctx.Export("FullEmailSyncURI", syncCF.Function.ServiceConfig.Uri())

	emailPushNotify, err := i.emailPushNotificationCF(gmailPubSub, syncCF)
	if err != nil {
		return err
	}

	// grant email push notification function invoke access to the sync function
	cloudrunv2.NewServiceIamMember(i.ctx, fmt.Sprintf("%s-can-invoke-%s", emailPushNotify.Name, syncCF.Name), &cloudrunv2.ServiceIamMemberArgs{
		Project:  pulumi.String(*i.Project.ProjectId),
		Location: pulumi.String(DefaultRegion),
		Name:     pulumi.String(syncCF.Name),
		Role:     pulumi.String("roles/run.invoker"),
		Member: emailPushNotify.ServiceAccount.Email.ApplyT(func(email string) (string, error) {
			return fmt.Sprintf("serviceAccount:%v", email), nil
		}).(pulumi.StringOutput),
	},
		pulumi.DependsOn([]pulumi.Resource{
			syncCF.Function,
			emailPushNotify.Function,
		}))

	return nil
}

func (i *Infra) fullEmailSyncCF() (*CloudFunction, error) {
	name := "full-email-sync"
	sa, err := i.createCloudFunctionServiceAccount(name)
	if err != nil {
		return nil, err
	}
	obj, err := i.uploadCloudFunction(name)
	if err != nil {
		return nil, err
	}

	cf, err := cloudfunctionsv2.NewFunction(i.ctx, name, &cloudfunctionsv2.FunctionArgs{
		Name: pulumi.String(name),
		// use the same location as the bucket
		Location:    pulumi.String(DefaultRegion),
		Project:     pulumi.String(*i.Project.ProjectId),
		Description: pulumi.String("Sync a user's historic emails starting from a given date"),
		BuildConfig: &cloudfunctionsv2.FunctionBuildConfigArgs{
			Runtime:    pulumi.String("go119"),
			EntryPoint: pulumi.String("FullEmailSync"),
			Source: &cloudfunctionsv2.FunctionBuildConfigSourceArgs{
				StorageSource: &cloudfunctionsv2.FunctionBuildConfigSourceStorageSourceArgs{
					Bucket: i.GCFBucket.Name,
					Object: obj.Name,
				},
			},
		},
		ServiceConfig: &cloudfunctionsv2.FunctionServiceConfigArgs{
			AvailableMemory:  pulumi.String("256M"),
			MinInstanceCount: pulumi.Int(0),
			MaxInstanceCount: pulumi.Int(1),
			TimeoutSeconds:   pulumi.Int(MaxHTTPTriggerTimeout),
			EnvironmentVariables: pulumi.StringMap{
				"SUPABASE_API_URL":           pulumi.String(i.config.Require("SUPABASE_API_URL")),
				"SUPABASE_API_KEY":           i.config.RequireSecret("SUPABASE_API_KEY"),
				"GOOGLE_OAUTH2_CREDENTIALS":  i.config.RequireSecret("GOOGLE_OAUTH2_CREDENTIALS"),
				"ML_SERVICE_URL":             i.config.RequireSecret("ML_SERVICE_URL"),
				"SENTRY_DSN":                 i.config.RequireSecret("SENTRY_DSN"),
				"EXAMPLES_GMAIL_OAUTH_TOKEN": i.config.RequireSecret("EXAMPLES_GMAIL_OAUTH_TOKEN"),
			},
			IngressSettings:            pulumi.String("ALLOW_ALL"),
			AllTrafficOnLatestRevision: pulumi.Bool(true),
			ServiceAccountEmail:        sa.Email,
		},
	}, pulumi.DependsOn([]pulumi.Resource{
		obj,
		sa,
	}))

	if err != nil {
		return nil, err
	}

	srv, err := cloudrun.LookupService(i.ctx, &cloudrun.LookupServiceArgs{
		Name:     name,
		Location: DefaultRegion,
		Project:  i.Project.ProjectId,
	})

	if err != nil {
		return nil, err
	}

	return &CloudFunction{
		Name:           name,
		ServiceAccount: sa,
		Function:       cf,
		Service:        srv,
	}, nil
}

func (i *Infra) emailPushNotificationCF(topic *pubsub.Topic, fullSync *CloudFunction) (*CloudFunction, error) {

	name := "email-push-notifications"
	sa, err := i.createCloudFunctionServiceAccount(name)
	if err != nil {
		return nil, err
	}
	obj, err := i.uploadCloudFunction(name)
	if err != nil {
		return nil, err
	}

	cf, err := cloudfunctionsv2.NewFunction(i.ctx, name, &cloudfunctionsv2.FunctionArgs{
		Name: pulumi.String(name),
		// use the same location as the bucket
		Location:    pulumi.String(DefaultRegion),
		Project:     pulumi.String(*i.Project.ProjectId),
		Description: pulumi.String("Handle user email push notifications"),
		BuildConfig: &cloudfunctionsv2.FunctionBuildConfigArgs{
			Runtime:    pulumi.String("go119"),
			EntryPoint: pulumi.String("EmailPushNotificationHandler"),
			Source: &cloudfunctionsv2.FunctionBuildConfigSourceArgs{
				StorageSource: &cloudfunctionsv2.FunctionBuildConfigSourceStorageSourceArgs{
					Bucket: i.GCFBucket.Name,
					Object: obj.Name,
				},
			},
		},
		ServiceConfig: &cloudfunctionsv2.FunctionServiceConfigArgs{
			AvailableMemory:  pulumi.String("256M"),
			MinInstanceCount: pulumi.Int(0),
			MaxInstanceCount: pulumi.Int(25),
			TimeoutSeconds:   pulumi.Int(MaxEventArcTriggerTimeout),
			EnvironmentVariables: pulumi.StringMap{
				"SUPABASE_API_URL":           pulumi.String(i.config.Require("SUPABASE_API_URL")),
				"SUPABASE_API_KEY":           i.config.RequireSecret("SUPABASE_API_KEY"),
				"GOOGLE_OAUTH2_CREDENTIALS":  i.config.RequireSecret("GOOGLE_OAUTH2_CREDENTIALS"),
				"ML_SERVICE_URL":             i.config.RequireSecret("ML_SERVICE_URL"),
				"SENTRY_DSN":                 i.config.RequireSecret("SENTRY_DSN"),
				"EXAMPLES_GMAIL_OAUTH_TOKEN": i.config.RequireSecret("EXAMPLES_GMAIL_OAUTH_TOKEN"),
				"TRIGGER_FULL_SYNC_URL":      fullSync.Function.ServiceConfig.Uri().Elem(),
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
		fullSync.Function,
	}))
	if err != nil {
		return nil, err
	}

	srv, err := cloudrun.LookupService(i.ctx, &cloudrun.LookupServiceArgs{
		Name:     name,
		Location: DefaultRegion,
		Project:  i.Project.ProjectId,
	})
	if err != nil {
		return nil, err
	}

	return &CloudFunction{
		Name:           name,
		ServiceAccount: sa,
		Function:       cf,
		Service:        srv,
	}, nil
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
