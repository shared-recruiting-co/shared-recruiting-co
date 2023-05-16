<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `label` directory is for defining and managing Gmail labels used by SRC for organizing and filtering emails. It includes candidate-specific labels for managing job opportunities and blocked emails, as well as labels for recruiters, such as `Recruiting` and `RecruitingOutbound`.

### Files
#### candidate.go
This file defines candidate-specific Gmail labels for managing job opportunities and blocked emails in the SRC recruiting platform. It includes a struct for the `CandidateLabels` and defines unique labels for Jobs, JobsOpportunity, JobsInterested, JobsNotInterested, JobsSaved, JobsVerified, and blocked emails.

#### label.go
This file defines Gmail labels used by SRC for organizing and filtering emails, including SRC, Allow, AllowSender, AllowDomain, Block, BlockSender, and BlockDomain. It imports the `gmail/v1` package and defines the `label` package.

#### recruiter.go
This file contains the `RecruiterLabels` struct, which includes all labels managed by SRC for recruiters. It also defines the `Recruiting` and `RecruitingOutbound` labels, which are used for all SRC-managed recruiting labels and identified recruiting outbound, respectively.

<!--- END SELFDOC --->