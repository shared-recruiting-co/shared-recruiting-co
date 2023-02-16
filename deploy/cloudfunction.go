package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudfunctionsv2"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudrun"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudrunv2"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudscheduler"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	serviceAccount "github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	DailyJobSchedule = "6 4 * * *"
)

type CloudFunction struct {
	Name           string
	ServiceAccount *serviceAccount.Account
	Function       *cloudfunctionsv2.Function
	Service        *cloudrun.LookupServiceResult
}

func (i *Infra) createCloudFunctions() error {
	syncCF, err := i.fullEmailSyncCF()
	if err != nil {
		return err
	}
	i.ctx.Export("FullEmailSyncURI", syncCF.Function.ServiceConfig.Uri())

	emailPushNotify, err := i.emailPushNotificationCF(syncCF)
	if err != nil {
		return err
	}

	// grant email push notification function invoke access to the sync function
	_, err = cloudrunv2.NewServiceIamMember(i.ctx, fmt.Sprintf("%s-can-invoke-%s", emailPushNotify.Name, syncCF.Name), &cloudrunv2.ServiceIamMemberArgs{
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

	if err != nil {
		return err
	}

	_, err = i.populateJobs()
	if err != nil {
		return err
	}

	_, err = i.watchCandidateEmails()
	if err != nil {
		return err
	}

	_, err = i.adhoc()
	if err != nil {
		return err
	}

	_, err = i.candidateGmailMessages()
	if err != nil {
		return err
	}

	_, err = i.recruiterGmailMessages()
	if err != nil {
		return err
	}

	return nil
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
			EntryPoint: pulumi.String("Handler"),
			EnvironmentVariables: pulumi.StringMap{
				// Use hash to force redeploy when code changes
				"FUNCTION_NAME":         pulumi.String(name),
				"FUNCTION_CONTENT_HASH": obj.Md5hash,
			},
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
				"SUPABASE_API_URL":          pulumi.String(i.config.Require("SUPABASE_API_URL")),
				"SUPABASE_API_KEY":          i.config.RequireSecret("SUPABASE_API_KEY"),
				"GOOGLE_OAUTH2_CREDENTIALS": i.config.RequireSecret("GOOGLE_OAUTH2_CREDENTIALS"),
				"SENTRY_DSN":                i.config.RequireSecret("SENTRY_DSN"),
				"GCP_PROJECT_ID":            pulumi.String(*i.Project.ProjectId),
				"CANDIDATE_GMAIL_MESSAGES_TOPIC": i.Topics.CandidateGmailMessages.Name.ApplyT(func(name string) string {
					return name
				}).(pulumi.StringOutput),
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

	_, err = pubsub.NewTopicIAMMember(i.ctx, fmt.Sprintf("%s-publish-to-candidate-gmail-messages", name), &pubsub.TopicIAMMemberArgs{
		Topic:   i.Topics.CandidateGmailMessages.ID(),
		Role:    pulumi.String("roles/pubsub.publisher"),
		Member:  pulumi.Sprintf("serviceAccount:%s", sa.Email),
		Project: pulumi.String(*i.Project.ProjectId),
	}, pulumi.DependsOn([]pulumi.Resource{
		cf,
		sa,
		i.Topics.CandidateGmailMessages,
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

func (i *Infra) emailPushNotificationCF(fullSync *CloudFunction) (*CloudFunction, error) {
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
			EntryPoint: pulumi.String("Handler"),
			EnvironmentVariables: pulumi.StringMap{
				// Use hash to force redeploy when code changes
				"FUNCTION_NAME":         pulumi.String(name),
				"FUNCTION_CONTENT_HASH": obj.Md5hash,
			},
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
				"SUPABASE_API_URL":          pulumi.String(i.config.Require("SUPABASE_API_URL")),
				"SUPABASE_API_KEY":          i.config.RequireSecret("SUPABASE_API_KEY"),
				"GOOGLE_OAUTH2_CREDENTIALS": i.config.RequireSecret("GOOGLE_OAUTH2_CREDENTIALS"),
				"SENTRY_DSN":                i.config.RequireSecret("SENTRY_DSN"),
				"TRIGGER_FULL_SYNC_URL":     fullSync.Function.ServiceConfig.Uri().Elem(),
				"GCP_PROJECT_ID":            pulumi.String(*i.Project.ProjectId),
				"CANDIDATE_GMAIL_MESSAGES_TOPIC": i.Topics.CandidateGmailMessages.Name.ApplyT(func(name string) string {
					return name
				}).(pulumi.StringOutput),
			},
			IngressSettings:            pulumi.String("ALLOW_INTERNAL_ONLY"),
			AllTrafficOnLatestRevision: pulumi.Bool(true),
			ServiceAccountEmail:        sa.Email,
		},
		EventTrigger: &cloudfunctionsv2.FunctionEventTriggerArgs{
			TriggerRegion: pulumi.String(DefaultRegion),
			PubsubTopic:   i.Topics.Gmail.ID(),
			EventType:     pulumi.String("google.cloud.pubsub.topic.v1.messagePublished"),
			// Disable retry
			RetryPolicy:         pulumi.String("RETRY_POLICY_DO_NOT_RETRY"),
			ServiceAccountEmail: sa.Email,
		},
	}, pulumi.DependsOn([]pulumi.Resource{
		i.Topics.Gmail,
		i.Topics.CandidateGmailMessages,
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

	_, err = pubsub.NewTopicIAMMember(i.ctx, fmt.Sprintf("%s-publish-to-candidate-gmail-messages", name), &pubsub.TopicIAMMemberArgs{
		Topic:   i.Topics.CandidateGmailMessages.ID(),
		Role:    pulumi.String("roles/pubsub.publisher"),
		Member:  pulumi.Sprintf("serviceAccount:%s", sa.Email),
		Project: pulumi.String(*i.Project.ProjectId),
	}, pulumi.DependsOn([]pulumi.Resource{
		cf,
		sa,
		i.Topics.CandidateGmailMessages,
	}))
	if err != nil {
		return nil, err
	}

	_, err = pubsub.NewTopicIAMMember(i.ctx, fmt.Sprintf("%s-publish-to-recruiter-gmail-messages", name), &pubsub.TopicIAMMemberArgs{
		Topic:   i.Topics.RecruiterGmailMessages.ID(),
		Role:    pulumi.String("roles/pubsub.publisher"),
		Member:  pulumi.Sprintf("serviceAccount:%s", sa.Email),
		Project: pulumi.String(*i.Project.ProjectId),
	}, pulumi.DependsOn([]pulumi.Resource{
		cf,
		sa,
		i.Topics.RecruiterGmailMessages,
	}))
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

func (i *Infra) candidateGmailMessages() (*CloudFunction, error) {
	name := "candidate-gmail-messages"
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
		Description: pulumi.String("Handle candidate gmail messages"),
		BuildConfig: &cloudfunctionsv2.FunctionBuildConfigArgs{
			Runtime:    pulumi.String("go119"),
			EntryPoint: pulumi.String("Handler"),
			EnvironmentVariables: pulumi.StringMap{
				// Use hash to force redeploy when code changes
				"FUNCTION_NAME":         pulumi.String(name),
				"FUNCTION_CONTENT_HASH": obj.Md5hash,
			},
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
			MaxInstanceCount: pulumi.Int(5),
			TimeoutSeconds:   pulumi.Int(MaxEventArcTriggerTimeout),
			EnvironmentVariables: pulumi.StringMap{
				"SUPABASE_API_URL":           pulumi.String(i.config.Require("SUPABASE_API_URL")),
				"SUPABASE_API_KEY":           i.config.RequireSecret("SUPABASE_API_KEY"),
				"GOOGLE_OAUTH2_CREDENTIALS":  i.config.RequireSecret("GOOGLE_OAUTH2_CREDENTIALS"),
				"ML_SERVICE_URL":             i.config.RequireSecret("ML_SERVICE_URL"),
				"SENTRY_DSN":                 i.config.RequireSecret("SENTRY_DSN"),
				"EXAMPLES_GMAIL_OAUTH_TOKEN": i.config.RequireSecret("EXAMPLES_GMAIL_OAUTH_TOKEN"),
			},
			IngressSettings:            pulumi.String("ALLOW_INTERNAL_ONLY"),
			AllTrafficOnLatestRevision: pulumi.Bool(true),
			ServiceAccountEmail:        sa.Email,
		},
		EventTrigger: &cloudfunctionsv2.FunctionEventTriggerArgs{
			TriggerRegion: pulumi.String(DefaultRegion),
			PubsubTopic:   i.Topics.CandidateGmailMessages.ID(),
			EventType:     pulumi.String("google.cloud.pubsub.topic.v1.messagePublished"),
			// Always retry failed messages
			RetryPolicy:         pulumi.String("RETRY_POLICY_RETRY"),
			ServiceAccountEmail: sa.Email,
		},
	}, pulumi.DependsOn([]pulumi.Resource{
		i.Topics.CandidateGmailMessages,
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

	// allow it to selectively re-push messages that fail
	_, err = pubsub.NewTopicIAMMember(i.ctx, fmt.Sprintf("%s-publish-to-candidate-gmail-messages", name), &pubsub.TopicIAMMemberArgs{
		Topic:   i.Topics.CandidateGmailMessages.ID(),
		Role:    pulumi.String("roles/pubsub.publisher"),
		Member:  pulumi.Sprintf("serviceAccount:%s", sa.Email),
		Project: pulumi.String(*i.Project.ProjectId),
	}, pulumi.DependsOn([]pulumi.Resource{
		cf,
		sa,
		i.Topics.CandidateGmailMessages,
	}))
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

func (i *Infra) recruiterGmailMessages() (*CloudFunction, error) {
	name := "recruiter-gmail-messages"
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
		Description: pulumi.String("Handle recruiter gmail messages"),
		BuildConfig: &cloudfunctionsv2.FunctionBuildConfigArgs{
			Runtime:    pulumi.String("go119"),
			EntryPoint: pulumi.String("Handler"),
			EnvironmentVariables: pulumi.StringMap{
				// Use hash to force redeploy when code changes
				"FUNCTION_NAME":         pulumi.String(name),
				"FUNCTION_CONTENT_HASH": obj.Md5hash,
			},
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
			MaxInstanceCount: pulumi.Int(5),
			TimeoutSeconds:   pulumi.Int(MaxEventArcTriggerTimeout),
			EnvironmentVariables: pulumi.StringMap{
				"SUPABASE_API_URL":          pulumi.String(i.config.Require("SUPABASE_API_URL")),
				"SUPABASE_API_KEY":          i.config.RequireSecret("SUPABASE_API_KEY"),
				"GOOGLE_OAUTH2_CREDENTIALS": i.config.RequireSecret("GOOGLE_OAUTH2_CREDENTIALS"),
				"ML_SERVICE_URL":            i.config.RequireSecret("ML_SERVICE_URL"),
				"SENTRY_DSN":                i.config.RequireSecret("SENTRY_DSN"),
			},
			IngressSettings:            pulumi.String("ALLOW_INTERNAL_ONLY"),
			AllTrafficOnLatestRevision: pulumi.Bool(true),
			ServiceAccountEmail:        sa.Email,
		},
		EventTrigger: &cloudfunctionsv2.FunctionEventTriggerArgs{
			TriggerRegion: pulumi.String(DefaultRegion),
			PubsubTopic:   i.Topics.RecruiterGmailMessages.ID(),
			EventType:     pulumi.String("google.cloud.pubsub.topic.v1.messagePublished"),
			// Always retry failed messages
			RetryPolicy:         pulumi.String("RETRY_POLICY_RETRY"),
			ServiceAccountEmail: sa.Email,
		},
	}, pulumi.DependsOn([]pulumi.Resource{
		i.Topics.RecruiterGmailMessages,
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

	// allow it to selectively re-push messages that fail
	_, err = pubsub.NewTopicIAMMember(i.ctx, fmt.Sprintf("%s-publish-to-recruiter-gmail-messages", name), &pubsub.TopicIAMMemberArgs{
		Topic:   i.Topics.RecruiterGmailMessages.ID(),
		Role:    pulumi.String("roles/pubsub.publisher"),
		Member:  pulumi.Sprintf("serviceAccount:%s", sa.Email),
		Project: pulumi.String(*i.Project.ProjectId),
	}, pulumi.DependsOn([]pulumi.Resource{
		cf,
		sa,
		i.Topics.RecruiterGmailMessages,
	}))
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

func (i *Infra) populateJobs() (*CloudFunction, error) {
	name := "populate-jobs"
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
		Description: pulumi.String("Parse inbound job emails and add them to the user's job board"),
		BuildConfig: &cloudfunctionsv2.FunctionBuildConfigArgs{
			Runtime:    pulumi.String("go119"),
			EntryPoint: pulumi.String("PopulateJobs"),
			EnvironmentVariables: pulumi.StringMap{
				// Use hash to force redeploy when code changes
				"FUNCTION_NAME":         pulumi.String(name),
				"FUNCTION_CONTENT_HASH": obj.Md5hash,
			},
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
				"SUPABASE_API_URL":          pulumi.String(i.config.Require("SUPABASE_API_URL")),
				"SUPABASE_API_KEY":          i.config.RequireSecret("SUPABASE_API_KEY"),
				"GOOGLE_OAUTH2_CREDENTIALS": i.config.RequireSecret("GOOGLE_OAUTH2_CREDENTIALS"),
				"ML_SERVICE_URL":            i.config.RequireSecret("ML_SERVICE_URL"),
				"SENTRY_DSN":                i.config.RequireSecret("SENTRY_DSN"),
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

	_, err = cloudscheduler.NewJob(i.ctx, fmt.Sprintf("%s-daily", name), &cloudscheduler.JobArgs{
		Name:        pulumi.Sprintf("%s-daily", name),
		Project:     pulumi.String(*i.Project.ProjectId),
		Region:      pulumi.String(DefaultRegion),
		Description: pulumi.String("Parse inbound job emails and add them to the user's job board each day"),
		HttpTarget: &cloudscheduler.JobHttpTargetArgs{
			Uri:        cf.ServiceConfig.Uri().Elem(),
			HttpMethod: pulumi.String(http.MethodPost),
			Headers: pulumi.StringMap{
				"Content-Type": pulumi.String("application/json"),
			},
			OidcToken: &cloudscheduler.JobHttpTargetOidcTokenArgs{
				ServiceAccountEmail: sa.Email,
				Audience:            cf.ServiceConfig.Uri().Elem(),
			},
		},
		Schedule: pulumi.String(DailyJobSchedule),
		// Set time zone to Greenwich Mean Time (GMT) aka UTC
		TimeZone:        pulumi.String("Etc/GMT"),
		AttemptDeadline: pulumi.String("1800s"),
	}, pulumi.DependsOn([]pulumi.Resource{
		cf,
	}))
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

func (i *Infra) watchCandidateEmails() (*CloudFunction, error) {
	name := "candidate-watch-emails"
	sa, err := i.createCloudFunctionServiceAccount(name)
	if err != nil {
		return nil, err
	}
	funcName := "watch-emails"
	obj, err := i.uploadCloudFunction(funcName)
	if err != nil {
		return nil, err
	}

	cf, err := cloudfunctionsv2.NewFunction(i.ctx, name, &cloudfunctionsv2.FunctionArgs{
		Name: pulumi.String(name),
		// use the same location as the bucket
		Location:    pulumi.String(DefaultRegion),
		Project:     pulumi.String(*i.Project.ProjectId),
		Description: pulumi.String("Subscribe to the candidate's email inbox and watch for relevant emails"),
		BuildConfig: &cloudfunctionsv2.FunctionBuildConfigArgs{
			Runtime:    pulumi.String("go119"),
			EntryPoint: pulumi.String("CandidateWatchEmails"),
			EnvironmentVariables: pulumi.StringMap{
				// Use hash to force redeploy when code changes
				"FUNCTION_NAME":         pulumi.String(name),
				"FUNCTION_CONTENT_HASH": obj.Md5hash,
			},
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
				"PUBSUB_TOPIC":              i.Topics.Gmail.ID(),
				"SUPABASE_API_URL":          pulumi.String(i.config.Require("SUPABASE_API_URL")),
				"SUPABASE_API_KEY":          i.config.RequireSecret("SUPABASE_API_KEY"),
				"GOOGLE_OAUTH2_CREDENTIALS": i.config.RequireSecret("GOOGLE_OAUTH2_CREDENTIALS"),
				"ML_SERVICE_URL":            i.config.RequireSecret("ML_SERVICE_URL"),
				"SENTRY_DSN":                i.config.RequireSecret("SENTRY_DSN"),
			},
			IngressSettings:            pulumi.String("ALLOW_ALL"),
			AllTrafficOnLatestRevision: pulumi.Bool(true),
			ServiceAccountEmail:        sa.Email,
		},
	}, pulumi.DependsOn([]pulumi.Resource{
		i.Topics.Gmail,
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

	_, err = cloudscheduler.NewJob(i.ctx, fmt.Sprintf("%s-daily", name), &cloudscheduler.JobArgs{
		Name:        pulumi.Sprintf("%s-daily", name),
		Project:     pulumi.String(*i.Project.ProjectId),
		Region:      pulumi.String(DefaultRegion),
		Description: pulumi.String("Subscribe daily to the user's email inbox and watch for new emails"),
		HttpTarget: &cloudscheduler.JobHttpTargetArgs{
			Uri:        cf.ServiceConfig.Uri().Elem(),
			HttpMethod: pulumi.String(http.MethodPost),
			Headers: pulumi.StringMap{
				"Content-Type": pulumi.String("application/json"),
			},
			OidcToken: &cloudscheduler.JobHttpTargetOidcTokenArgs{
				ServiceAccountEmail: sa.Email,
				Audience:            cf.ServiceConfig.Uri().Elem(),
			},
		},
		Schedule: pulumi.String(DailyJobSchedule),
		// Set time zone to Greenwich Mean Time (GMT) aka UTC
		TimeZone:        pulumi.String("Etc/GMT"),
		AttemptDeadline: pulumi.String("1800s"),
	}, pulumi.DependsOn([]pulumi.Resource{
		cf,
	}))
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

func (i *Infra) adhoc() (*CloudFunction, error) {
	name := "adhoc"
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
		Description: pulumi.String("Adhoc function to run a one-off, manual tasks"),
		BuildConfig: &cloudfunctionsv2.FunctionBuildConfigArgs{
			Runtime:    pulumi.String("go119"),
			EntryPoint: pulumi.String("Reclassify"),
			EnvironmentVariables: pulumi.StringMap{
				// Use hash to force redeploy when code changes
				"FUNCTION_NAME":         pulumi.String(name),
				"FUNCTION_CONTENT_HASH": obj.Md5hash,
			},
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
