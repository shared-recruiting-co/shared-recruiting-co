# Auto-Generated Documentation 

## Summary

The `deploy` directory contains files related to the deployment of the SRC recruiting platform. It includes configuration files for Pulumi, a cloud infrastructure management tool, and GolangCI-Lint, a Go linter. The directory also contains Go files that create and configure cloud functions for various tasks in the platform, as well as unit tests for some of these functions. Additionally, there are files that set up a storage bucket and topics for candidate and recruiter Gmail subscriptions using the Google Cloud Pub/Sub service. The `main.go` file contains the main function that sets up the infrastructure for the platform.


## Technologies 

The technologies used by the `deploy` directory are:

- Pulumi: a tool for infrastructure as code that is used to define the infrastructure for SRC's GCP deployment and to specify the production deployment configuration, including secure credentials for cloud services, secrets provider, encrypted key, and cloud function configuration.
- Google Cloud Platform: the cloud platform used to deploy the SRC recruiting platform.
- Go: the programming language used to write the code for creating and configuring cloud functions and setting up topics for Gmail subscriptions and messages.
- GolangCI-Lint: a tool for linting Go code that is used to exclude certain issues from the linting process.
- Google Cloud Pub/Sub: a service used to set up topics for candidate and recruiter Gmail subscriptions, messages, and label changes.

## Files

#### Pulumi.prod.yaml 
This file contains the production deployment configuration for the SRC recruiting platform, including secure credentials for cloud services, secrets provider, encrypted key, and cloud function configuration.

#### Pulumi.yaml
This file is used to define the infrastructure for SRC's GCP deployment. It specifies the name of the cloud functions and the runtime environment as Go.

#### .golangci.yml
The `.golangci.yml` file contains configuration settings for the GolangCI-Lint tool. The `exclude` section lists regular expressions of issue texts to exclude from the linting process. In this case, the `loop variable tc captured by func literal` issue is excluded.

#### cloudfunction.go
This file contains functions that create and configure cloud functions for various tasks in the SRC recruiting platform, including syncing emails, handling push notifications, and subscribing to email inboxes. The functions create a cloud function with a unique name, location, and description, and set various environment variables and configurations such as memory, instance count, and timeout. They also create a service account and upload the cloud function's source code to a storage bucket.

#### cloudfunction_test.go
The `cloudfunction_test.go` file contains unit tests for the `shortenAccountId` function. The function takes a string input and returns a shortened version of it based on certain rules. The tests cover different scenarios and expected outputs for the function. The `testing` package is imported to run the tests.

#### go.mod
This file specifies the Go module for the Shared Recruiting Co. deployment directory. It lists the required dependencies and their versions, including the Pulumi and Pulumi GCP SDKs.

#### main.go
The `main.go` file contains the main function that sets up the infrastructure for the Shared Recruiting Co. platform. It imports various packages from the `pulumi-gcp/sdk/v6/go/gcp` and `pulumi/sdk/v3/go` libraries. The `main` function creates a new Pulumi stack and sets up a storage bucket, grants publish permission to a Gmail topic, and creates cloud functions. The `Infra` struct contains the context, configuration, project, GCF bucket, and topics. The `Topics` struct contains various Gmail topics used by the platform. The `MaxEventArcTriggerTimeout` and `MaxHTTPTriggerTimeout` variables define the maximum timeout values for event arc and HTTP triggers, respectively. The `DefaultRegion` variable defines the default region for the storage bucket. The `TODOs` section lists tasks that need to be completed.

#### pubsub.go
The `pubsub.go` file contains a function called `setupTopics()` that sets up topics for candidate and recruiter Gmail subscriptions, messages, and label changes using the Google Cloud Pub/Sub service. The function creates new topics with specified names and project IDs, and enforces schema validation. The function returns an error if any of the topic creations fail.

