package gmail

import (
	"google.golang.org/api/gmail/v1"
)

const (
	SRC_Label                    = "@SRC"
	SRC_JobOpportunityLabel      = "@SRC/Job Opportunity"
	SRC_Color                    = "#ff7537"
	SRC_JobOpportunityLabelColor = "#16a765"
	white                        = "#ffffff"
)

func GetOrCreateLabel(srv *gmail.Service, userID string, labelID string, backgroundColor string, textColor string) (*gmail.Label, error) {
	labels, err := srv.Users.Labels.List(userID).Do()

	if err != nil {
		return nil, err
	}

	for _, label := range labels.Labels {
		if label.Name == labelID {
			return label, err
		}
	}

	label, err := srv.Users.Labels.Create(userID, &gmail.Label{Name: labelID, Color: &gmail.LabelColor{
		BackgroundColor: backgroundColor,
		TextColor:       textColor,
	}}).Do()

	if err != nil {
		return nil, err
	}

	return label, err
}

func GetOrCreateSRCLabel(srv *gmail.Service, userID string) (*gmail.Label, error) {
	return GetOrCreateLabel(srv, userID, SRC_Label, SRC_Color, white)
}

func GetOrCreateSRCJobOpportunityLabel(srv *gmail.Service, userID string) (*gmail.Label, error) {
	return GetOrCreateLabel(srv, userID, SRC_JobOpportunityLabel, SRC_JobOpportunityLabelColor, white)
}
