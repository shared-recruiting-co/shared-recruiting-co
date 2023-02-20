package label

import "google.golang.org/api/gmail/v1"

// RecruiterLabels contains all labels managed by @SRC for recruiters..
// IDs are unique per gmail account.
// The rest of the fields should be the same across all accounts.
type RecruiterLabels struct {
	// SRC is the @SRC parent label for all @SRC labels.
	SRC *gmail.Label
	// Recruiting is the @SRC/Recruiting folder for all @SRC job labels.
	Recruiting *gmail.Label
	// RecruitingOutbound is the @SRC/Recruiting/Outbound label for identified recruiting outbound.
	RecruitingOutbound *gmail.Label
}

// Recruiter-specific labels
var (
	// Recruiting is the @SRC/Recruiting folder for all @SRC managed recruiting labels.
	Recruiting = gmail.Label{
		Name: SRC.Name + "/Recruiting",
		// Hide folder labels
		MessageListVisibility: "hide",
		Color:                 SRC.Color,
	}
	// RecruitingOutbound is the @SRC/Jobs/Opportunity label for identified recruiting outbound..
	RecruitingOutbound = gmail.Label{
		Name: Recruiting.Name + "/Outbound",
		// Show leaf labels
		MessageListVisibility: "show",
		// use same color as parent
		Color: Recruiting.Color,
	}
)
