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
	candidateEmailSync, err := i.candidateEmailSyncCF()
	if err != nil {
		return err
	}

	recruiterEmailSync, err := i.recruiterEmailSyncCF()
	if err != nil {
		return err
	}

	_, err = i.candidateGmailPushNotifications(candidateEmailSync)
	if err != nil {
		return err
	}

	_, err = i.recruiterGmailPushNotifications(recruiterEmailSync)
	if err != nil {
		return err
	}

	_, err = i.candidateGmailSubscription()
	if err != nil {
		return err
	}

	_, err = i.recruiterGmailSubscription()
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

	_, err = i.scrapeFromYcJobs()
	if err != nil {
		return err
	}

	_, err = i.candidateGmailLabelChanges()
	if err != nil {
		return err
	}

	return nil
}

func (i *Infra) uploadCloudFunction(funcName, objName string) (*storage.BucketObject, error) {
	if objName == "" {
		objName = funcName
	}
	snakeName := strings.ReplaceAll(funcName, "-", "_")
	// for now hardcode path
	dir := "../cloudfunctions"

	return storage.NewBucketObject(i.ctx, fmt.Sprintf("cf-upload-%s", objName), &storage.BucketObjectArgs{
		Name:   pulumi.Sprintf("%s/function.zip", objName),
		Bucket: i.GCFBucket.Name,
		Source: pulumi.NewFileArchive(fmt.Sprintf("%s/%s", dir, snakeName)),
	}, pulumi.DependsOn([]pulumi.Resource{
		i.GCFBucket,
	}))
}

func shortenAccountId(id string) string {
	// replace common words with abbreviations
	id = strings.ReplaceAll(id, "gmail", "gm")
	id = strings.ReplaceAll(id, "candidate", "ca")
	id = strings.ReplaceAll(id, "recruiter", "re")

	if len(id) < 30 {
		return id
	}

	return id[:30]
}

func (i *Infra) scheduleDailyCloudFunction(name, desc string, sa *serviceAccount.Account, cf *cloudfunctionsv2.Function) error {
	_, err := cloudscheduler.NewJob(i.ctx, fmt.Sprintf("%s-daily", name), &cloudscheduler.JobArgs{
		Name:        pulumi.Sprintf("%s-daily", name),
		Project:     pulumi.String(*i.Project.ProjectId),
		Region:      pulumi.String(DefaultRegion),
		Description: pulumi.String(desc),
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

	return err
}

func (i *Infra) createCloudFunctionServiceAccount(name string) (*serviceAccount.Account, error) {
	account := fmt.Sprintf("sa-cf-%s", name)
	accountId := account

	if len(accountId) > 30 {
		accountId = shortenAccountId(accountId)
	}

	sa, err := serviceAccount.NewAccount(i.ctx, account, &serviceAccount.AccountArgs{
		Project: pulumi.String(*i.Project.ProjectId),
		// AccountId be 6-30 characters long and match the regular expression [a-z]([-a-z0-9]*[a-z0-9])?
		// https://cloud.google.com/iam/docs/service-accounts#creating_a_service_account
		AccountId:   pulumi.String(accountId),
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

func (i *Infra) scrapeFromYcJobs() (*CloudFunction, error) {
	name := "scrape-job-listings"
	sa, err := i.createCloudFunctionServiceAccount(name)
	if err != nil {
		return nil, err
	}
	obj, err := i.uploadCloudFunction(name, "")
	if err != nil {
		return nil, err
	}

	cf, err := cloudfunctionsv2.NewFunction(i.ctx, name, &cloudfunctionsv2.FunctionArgs{
		Name: pulumi.String(name),
		// use the same location as the bucket
		Location:    pulumi.String(DefaultRegion),
		Project:     pulumi.String(*i.Project.ProjectId),
		Description: pulumi.String("Scrape job listings from YC"),
		BuildConfig: &cloudfunctionsv2.FunctionBuildConfigArgs{
			Runtime:    pulumi.String("go120"),
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
			TimeoutSeconds:   pulumi.Int(MaxEventArcTriggerTimeout),
			EnvironmentVariables: pulumi.StringMap{
				"GOOGLE_OAUTH2_CREDENTIALS": i.config.RequireSecret("GOOGLE_OAUTH2_CREDENTIALS"),
				"SENTRY_DSN":                i.config.RequireSecret("SENTRY_DSN"),
			},
			IngressSettings:            pulumi.String("ALLOW_INTERNAL_ONLY"),
			AllTrafficOnLatestRevision: pulumi.Bool(true),
			ServiceAccountEmail:        sa.Email,
		},
		EventTrigger: &cloudfunctionsv2.FunctionEventTriggerArgs{
			TriggerRegion: pulumi.String(DefaultRegion),
			PubsubTopic:   i.Topics.ScrapeJobListings.ID(),
			EventType:     pulumi.String("google.cloud.pubsub.topic.v1.messagePublished"),
			// Always retry failed messages
			RetryPolicy:         pulumi.String("RETRY_POLICY_RETRY"),
			ServiceAccountEmail: sa.Email,
		},
	}, pulumi.DependsOn([]pulumi.Resource{
		i.Topics.ScrapeJobListings,
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

	err = i.scheduleDailyCloudFunction(name, "Daily scrape new jobs from YC", sa, cf)
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

func (i *Infra) candidateEmailSyncCF() (*CloudFunction, error) {
	name := "candidate-email-sync"
	sa, err := i.createCloudFunctionServiceAccount(name)
	if err != nil {
		return nil, err
	}
	obj, err := i.uploadCloudFunction(name, "")
	if err != nil {
		return nil, err
	}

	cf, err := cloudfunctionsv2.NewFunction(i.ctx, name, &cloudfunctionsv2.FunctionArgs{
		Name: pulumi.String(name),
		// use the same location as the bucket
		Location:    pulumi.String(DefaultRegion),
		Project:     pulumi.String(*i.Project.ProjectId),
		Description: pulumi.String("Sync a candidate's historic emails starting from a given date"),
		BuildConfig: &cloudfunctionsv2.FunctionBuildConfigArgs{
			Runtime:    pulumi.String("go120"),
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
				"GMAIL_MESSAGES_TOPIC": i.Topics.CandidateGmailMessages.Name.ApplyT(func(name string) string {
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

func (i *Infra) recruiterEmailSyncCF() (*CloudFunction, error) {
	name := "recruiter-email-sync"
	sa, err := i.createCloudFunctionServiceAccount(name)
	if err != nil {
		return nil, err
	}
	obj, err := i.uploadCloudFunction(name, "")
	if err != nil {
		return nil, err
	}

	cf, err := cloudfunctionsv2.NewFunction(i.ctx, name, &cloudfunctionsv2.FunctionArgs{
		Name: pulumi.String(name),
		// use the same location as the bucket
		Location:    pulumi.String(DefaultRegion),
		Project:     pulumi.String(*i.Project.ProjectId),
		Description: pulumi.String("Sync a recruiter's historic emails starting from a given date"),
		BuildConfig: &cloudfunctionsv2.FunctionBuildConfigArgs{
			Runtime:    pulumi.String("go120"),
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
				"GMAIL_MESSAGES_TOPIC": i.Topics.RecruiterGmailMessages.Name.ApplyT(func(name string) string {
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

func (i *Infra) candidateGmailPushNotifications(emailSync *CloudFunction) (*CloudFunction, error) {
	name := "candidate-gmail-push-notifications"
	sa, err := i.createCloudFunctionServiceAccount(name)
	if err != nil {
		return nil, err
	}
	obj, err := i.uploadCloudFunction(name, "")
	if err != nil {
		return nil, err
	}

	cf, err := cloudfunctionsv2.NewFunction(i.ctx, name, &cloudfunctionsv2.FunctionArgs{
		Name: pulumi.String(name),
		// use the same location as the bucket
		Location:    pulumi.String(DefaultRegion),
		Project:     pulumi.String(*i.Project.ProjectId),
		Description: pulumi.String("Handle candidate gmail push notifications"),
		BuildConfig: &cloudfunctionsv2.FunctionBuildConfigArgs{
			Runtime:    pulumi.String("go120"),
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
				"TRIGGER_EMAIL_SYNC_URL":    emailSync.Function.ServiceConfig.Uri().Elem(),
				"GCP_PROJECT_ID":            pulumi.String(*i.Project.ProjectId),
				"GMAIL_MESSAGES_TOPIC": i.Topics.CandidateGmailMessages.Name.ApplyT(func(name string) string {
					return name
				}).(pulumi.StringOutput),
				"GMAIL_LABEL_CHANGES_TOPIC": i.Topics.CandidateGmailLabelChanges.Name.ApplyT(func(name string) string {
					return name
				}).(pulumi.StringOutput),
			},
			IngressSettings:            pulumi.String("ALLOW_INTERNAL_ONLY"),
			AllTrafficOnLatestRevision: pulumi.Bool(true),
			ServiceAccountEmail:        sa.Email,
		},
		EventTrigger: &cloudfunctionsv2.FunctionEventTriggerArgs{
			TriggerRegion: pulumi.String(DefaultRegion),
			PubsubTopic:   i.Topics.CandidateGmailSubscription.ID(),
			EventType:     pulumi.String("google.cloud.pubsub.topic.v1.messagePublished"),
			// Disable retry
			RetryPolicy:         pulumi.String("RETRY_POLICY_DO_NOT_RETRY"),
			ServiceAccountEmail: sa.Email,
		},
	}, pulumi.DependsOn([]pulumi.Resource{
		i.Topics.CandidateGmailSubscription,
		i.Topics.CandidateGmailMessages,
		i.Topics.CandidateGmailLabelChanges,
		obj,
		sa,
		emailSync.Function,
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

	// grant function invoke access to the gmail sync function
	_, err = cloudrunv2.NewServiceIamMember(i.ctx, fmt.Sprintf("%s-can-invoke-%s", name, emailSync.Name), &cloudrunv2.ServiceIamMemberArgs{
		Project:  pulumi.String(*i.Project.ProjectId),
		Location: pulumi.String(DefaultRegion),
		Name:     pulumi.String(emailSync.Name),
		Role:     pulumi.String("roles/run.invoker"),
		Member: sa.Email.ApplyT(func(email string) (string, error) {
			return fmt.Sprintf("serviceAccount:%v", email), nil
		}).(pulumi.StringOutput),
	},
		pulumi.DependsOn([]pulumi.Resource{
			emailSync.Function,
			cf,
		}))
	if err != nil {
		return nil, err
	}

	// Grant publish permission to the necessary topics
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

	_, err = pubsub.NewTopicIAMMember(i.ctx, fmt.Sprintf("%s-publish-to-candidate-gmail-label-changes", name), &pubsub.TopicIAMMemberArgs{
		Topic:   i.Topics.CandidateGmailLabelChanges.ID(),
		Role:    pulumi.String("roles/pubsub.publisher"),
		Member:  pulumi.Sprintf("serviceAccount:%s", sa.Email),
		Project: pulumi.String(*i.Project.ProjectId),
	}, pulumi.DependsOn([]pulumi.Resource{
		cf,
		sa,
		i.Topics.CandidateGmailLabelChanges,
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
	obj, err := i.uploadCloudFunction(name, "")
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
			Runtime:    pulumi.String("go120"),
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
	obj, err := i.uploadCloudFunction(name, "")
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
			Runtime:    pulumi.String("go120"),
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
			// Right now, template detection is sensitive to ordering of messages
			// so we need to keep this at 1
			// Long term, we should be able to scale this up
			MaxInstanceCount: pulumi.Int(1),
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

func (i *Infra) recruiterGmailPushNotifications(emailSync *CloudFunction) (*CloudFunction, error) {
	name := "recruiter-gmail-push-notifications"
	sa, err := i.createCloudFunctionServiceAccount(name)
	if err != nil {
		return nil, err
	}
	obj, err := i.uploadCloudFunction(name, "")
	if err != nil {
		return nil, err
	}

	cf, err := cloudfunctionsv2.NewFunction(i.ctx, name, &cloudfunctionsv2.FunctionArgs{
		Name: pulumi.String(name),
		// use the same location as the bucket
		Location:    pulumi.String(DefaultRegion),
		Project:     pulumi.String(*i.Project.ProjectId),
		Description: pulumi.String("Handle recruiter gmail push notifications"),
		BuildConfig: &cloudfunctionsv2.FunctionBuildConfigArgs{
			Runtime:    pulumi.String("go120"),
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
				"TRIGGER_EMAIL_SYNC_URL":    emailSync.Function.ServiceConfig.Uri().Elem(),
				"GCP_PROJECT_ID":            pulumi.String(*i.Project.ProjectId),
				"GMAIL_MESSAGES_TOPIC": i.Topics.RecruiterGmailMessages.Name.ApplyT(func(name string) string {
					return name
				}).(pulumi.StringOutput),
			},
			IngressSettings:            pulumi.String("ALLOW_INTERNAL_ONLY"),
			AllTrafficOnLatestRevision: pulumi.Bool(true),
			ServiceAccountEmail:        sa.Email,
		},
		EventTrigger: &cloudfunctionsv2.FunctionEventTriggerArgs{
			TriggerRegion: pulumi.String(DefaultRegion),
			PubsubTopic:   i.Topics.RecruiterGmailSubscription.ID(),
			EventType:     pulumi.String("google.cloud.pubsub.topic.v1.messagePublished"),
			// Disable retry
			RetryPolicy:         pulumi.String("RETRY_POLICY_DO_NOT_RETRY"),
			ServiceAccountEmail: sa.Email,
		},
	}, pulumi.DependsOn([]pulumi.Resource{
		i.Topics.RecruiterGmailSubscription,
		i.Topics.RecruiterGmailMessages,
		obj,
		sa,
		emailSync.Function,
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

	// grant function invoke access to the gmail sync function
	_, err = cloudrunv2.NewServiceIamMember(i.ctx, fmt.Sprintf("%s-can-invoke-%s", name, emailSync.Name), &cloudrunv2.ServiceIamMemberArgs{
		Project:  pulumi.String(*i.Project.ProjectId),
		Location: pulumi.String(DefaultRegion),
		Name:     pulumi.String(emailSync.Name),
		Role:     pulumi.String("roles/run.invoker"),
		Member: sa.Email.ApplyT(func(email string) (string, error) {
			return fmt.Sprintf("serviceAccount:%v", email), nil
		}).(pulumi.StringOutput),
	},
		pulumi.DependsOn([]pulumi.Resource{
			emailSync.Function,
			cf,
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

func (i *Infra) candidateGmailSubscription() (*CloudFunction, error) {
	name := "candidate-gmail-subscription"
	sa, err := i.createCloudFunctionServiceAccount(name)
	if err != nil {
		return nil, err
	}
	funcName := "gmail-subscription"
	obj, err := i.uploadCloudFunction(funcName, name)
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
			Runtime:    pulumi.String("go120"),
			EntryPoint: pulumi.String("CandidateGmailSubscription"),
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
				"PUBSUB_TOPIC":              i.Topics.CandidateGmailSubscription.ID(),
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
		i.Topics.CandidateGmailSubscription,
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

	err = i.scheduleDailyCloudFunction(name, "Subscribe daily to the user's email inbox and watch for new emails", sa, cf)
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

func (i *Infra) recruiterGmailSubscription() (*CloudFunction, error) {
	name := "recruiter-gmail-subscription"
	sa, err := i.createCloudFunctionServiceAccount(name)
	if err != nil {
		return nil, err
	}
	funcName := "gmail-subscription"
	obj, err := i.uploadCloudFunction(funcName, name)
	if err != nil {
		return nil, err
	}

	cf, err := cloudfunctionsv2.NewFunction(i.ctx, name, &cloudfunctionsv2.FunctionArgs{
		Name: pulumi.String(name),
		// use the same location as the bucket
		Location:    pulumi.String(DefaultRegion),
		Project:     pulumi.String(*i.Project.ProjectId),
		Description: pulumi.String("Subscribe to the recruiter's email inbox and watch for relevant emails"),
		BuildConfig: &cloudfunctionsv2.FunctionBuildConfigArgs{
			Runtime:    pulumi.String("go120"),
			EntryPoint: pulumi.String("RecruiterGmailSubscription"),
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
				"PUBSUB_TOPIC":              i.Topics.RecruiterGmailSubscription.ID(),
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
		i.Topics.RecruiterGmailSubscription,
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

	err = i.scheduleDailyCloudFunction(name, "Subscribe daily to the user's email inbox and watch for new emails", sa, cf)
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
	obj, err := i.uploadCloudFunction(name, "")
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
			Runtime:    pulumi.String("go120"),
			EntryPoint: pulumi.String("Unsubscribe"),
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

func (i *Infra) candidateGmailLabelChanges() (*CloudFunction, error) {
	name := "candidate-gmail-label-changes"
	sa, err := i.createCloudFunctionServiceAccount(name)
	if err != nil {
		return nil, err
	}
	obj, err := i.uploadCloudFunction(name, "")
	if err != nil {
		return nil, err
	}

	cf, err := cloudfunctionsv2.NewFunction(i.ctx, name, &cloudfunctionsv2.FunctionArgs{
		Name: pulumi.String(name),
		// use the same location as the bucket
		Location:    pulumi.String(DefaultRegion),
		Project:     pulumi.String(*i.Project.ProjectId),
		Description: pulumi.String("Handle candidate gmail label changes"),
		BuildConfig: &cloudfunctionsv2.FunctionBuildConfigArgs{
			Runtime:    pulumi.String("go120"),
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
			PubsubTopic:   i.Topics.CandidateGmailLabelChanges.ID(),
			EventType:     pulumi.String("google.cloud.pubsub.topic.v1.messagePublished"),
			// Always retry failed messages
			RetryPolicy:         pulumi.String("RETRY_POLICY_RETRY"),
			ServiceAccountEmail: sa.Email,
		},
	}, pulumi.DependsOn([]pulumi.Resource{
		i.Topics.CandidateGmailLabelChanges,
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
