package gmail

import (
	"google.golang.org/api/gmail/v1"
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
		// Always user white text for now...
		TextColor: textColor,
	}}).Do()

	if err != nil {
		return nil, err
	}

	return label, err
}
