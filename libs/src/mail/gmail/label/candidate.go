package label

import (
	"google.golang.org/api/gmail/v1"
)

// CandidateLabels contains all labels managed by @SRC for candidates.
// IDs are unique per gmail account.
// The rest of the fields should be the same across all accounts.
type CandidateLabels struct {
	// SRC is the @SRC parent label for all @SRC labels.
	SRC *gmail.Label
	// Jobs is the @SRC/Jobs folder for all @SRC job labels.
	Jobs *gmail.Label
	// JobsOpportunity is the @SRC/Jobs/Opportunity label for identified job opportunities.
	JobsOpportunity *gmail.Label
	// JobsInterested is the @SRC/Jobs/Interested label for job opportunities the candidate marked as interested.
	JobsInterested *gmail.Label
	// JobsNotInterested is the @SRC/Jobs/Not Interested label for job opportunities the candidate marked as not interested.
	JobsNotInterested *gmail.Label
	// JobsSaved is the @SRC/Jobs/Saved label for job opportunities the candidate saved.
	JobsSaved *gmail.Label
	// JobsVerified is the @SRC/Jobs/Verified label for job opportunities that are sent from recruiters using SRC
	JobsVerified *gmail.Label
	// Allow is the @SRC/Allow folder for senders and domains that always bypass SRC
	Allow *gmail.Label
	// AllowSender is the @SRC/Allow/Sender folder for senders that always bypass SRC
	AllowSender *gmail.Label
	// AllowDomain is the @SRC/Allow/Domain folder for email domains that always bypass SRC
	AllowDomain *gmail.Label
	// Block is the @SRC/Allow folder for senders and domains that are automatically removed from your inbox.
	// Blocked emails are not analyzed by SRC.
	Block *gmail.Label
	// BlockSender is the @SRC/Block/Sender folder for senders that are automatically removed from your inbox.
	// Blocked emails are not analyzed by SRC.
	BlockSender *gmail.Label
	// BlockDomain is the @SRC/Block/Domain folder for email domains that are automatically removed from your inbox.
	// Blocked emails are not analyzed by SRC.
	BlockDomain *gmail.Label
	// BlockGraveyard is the @SRC/Block/🪦 folder for emails that have been removed from your inbox.
	// Blocked emails are not be analyzed by SRC.
	BlockGraveyard *gmail.Label
}

// Candidate-specific labels
var (
	// Jobs is the @SRC/Jobs folder for all @SRC job labels.
	Jobs = gmail.Label{
		Name: SRC.Name + "/Jobs",
		// Hide folder labels
		MessageListVisibility: "hide",
		Color:                 SRC.Color,
	}
	// JobsOpportunity is the @SRC/Jobs/Opportunity label for identified job opportunities.
	JobsOpportunity = gmail.Label{
		Name: Jobs.Name + "/Opportunity",
		// Show leaf labels
		MessageListVisibility: "show",
		// use same color as parent
		Color: Jobs.Color,
	}
	// JobsInterested is the @SRC/Jobs/Interested label for job opportunities the candidate marked as interested.
	JobsInterested = gmail.Label{
		Name: Jobs.Name + "/Interested",
		// Show leaf labels
		MessageListVisibility: "show",
		// green
		Color: &gmail.LabelColor{
			BackgroundColor: "#16a765",
			TextColor:       "#ffffff",
		},
	}
	// JobsNotInterested is the @SRC/Jobs/Not Interested label for job opportunities the candidate marked as not interested.
	JobsNotInterested = gmail.Label{
		Name: Jobs.Name + "/Not Interested",
		// Show leaf labels
		MessageListVisibility: "show",
		// red
		Color: &gmail.LabelColor{
			BackgroundColor: "#cc3a21",
			TextColor:       "#ffffff",
		},
	}
	// JobsSaved is the @SRC/Jobs/Saved label for job opportunities the candidate saved.
	JobsSaved = gmail.Label{
		Name: Jobs.Name + "/Saved",
		// Show leaf labels
		MessageListVisibility: "show",
		// yellow
		Color: &gmail.LabelColor{
			BackgroundColor: "#ffad46",
			TextColor:       "#ffffff",
		},
	}
	// JobsVerified is the @SRC/Jobs/Verified label for verified job opportunities.
	JobsVerified = gmail.Label{
		Name: Jobs.Name + "/Verified",
		// Show leaf labels
		MessageListVisibility: "show",
		// gray-ish black with white text
		Color: &gmail.LabelColor{
			BackgroundColor: "#333633",
			TextColor:       "#ffffff",
		},
	}
)
