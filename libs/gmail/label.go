package gmail

import (
	"google.golang.org/api/gmail/v1"
)

const (
	SRCLabel               = "@SRC"
	SRCJobOpportunityLabel = "@SRC/Job Opportunity"
)

var (
	SRCLabelColor = &gmail.LabelColor{
		BackgroundColor: "#4986e7",
		TextColor:       "#ffffff",
	}
	SRCJobOpportunityLabelColor = &gmail.LabelColor{
		BackgroundColor: "#c9daf8",
		TextColor:       "#1c4587",
	}
)
