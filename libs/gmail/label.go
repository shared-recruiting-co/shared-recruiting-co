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

func GetOrCreateLabel(srv *gmail.Service, userID string, labelID string, color *gmail.LabelColor) (*gmail.Label, error) {
	labels, err := srv.Users.Labels.List(userID).Do()

	if err != nil {
		return nil, err
	}

	for _, label := range labels.Labels {
		if label.Name == labelID {
			return label, err
		}
	}

	label, err := srv.Users.Labels.Create(userID, &gmail.Label{Name: labelID, Color: color}).Do()

	if err != nil {
		return nil, err
	}

	return label, err
}

func GetOrCreateSRCLabel(srv *gmail.Service, userID string) (*gmail.Label, error) {
	return GetOrCreateLabel(srv, userID, SRCLabel, SRCLabelColor)
}

func GetOrCreateSRCJobOpportunityLabel(srv *gmail.Service, userID string) (*gmail.Label, error) {
	return GetOrCreateLabel(srv, userID, SRCJobOpportunityLabel, SRCJobOpportunityLabelColor)
}
