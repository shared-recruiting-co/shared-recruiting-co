package label

import (
	"google.golang.org/api/gmail/v1"
)

// Generic labels
var (
	// SRC is the @SRC parent label for all @SRC labels.
	SRC = gmail.Label{
		Name:                "@SRC",
		LabelListVisibility: "labelShow",
		// Hide folder labels
		MessageListVisibility: "hide",
		Color: &gmail.LabelColor{
			BackgroundColor: "#4986e7",
			TextColor:       "#ffffff",
		},
	}
	// Allow is the @SRC/Allow folder for senders and domains that always bypass SRC
	Allow = gmail.Label{
		Name: SRC.Name + "/Allow",
		// Hide folder labels
		MessageListVisibility: "hide",
		Color: &gmail.LabelColor{
			BackgroundColor: "#16a765",
			TextColor:       "#ffffff",
		},
	}
	// AllowSender is the @SRC/Allow/Sender folder for senders that always bypass SRC
	AllowSender = gmail.Label{
		Name: Allow.Name + "/Sender",
		// Show leaf labels
		MessageListVisibility: "show",
		Color:                 Allow.Color,
	}
	// AllowDomain is the @SRC/Allow/Domain folder for email domains that always bypass SRC
	AllowDomain = gmail.Label{
		Name: Allow.Name + "/Domain",
		// Show leaf labels
		MessageListVisibility: "show",
		Color:                 Allow.Color,
	}
	// Block is the @SRC/Allow folder for senders and domains that are automatically removed from your inbox.
	// Blocked emails are not analyzed by SRC.
	Block = gmail.Label{
		Name: SRC.Name + "/Block",
		// Hide folder labels
		MessageListVisibility: "hide",
		Color: &gmail.LabelColor{
			BackgroundColor: "#cc3a21",
			TextColor:       "#ffffff",
		},
	}
	// BlockSender is the @SRC/Block/Sender folder for senders that are automatically removed from your inbox.
	// Blocked emails are not analyzed by SRC.
	BlockSender = gmail.Label{
		Name:                  Block.Name + "/Sender",
		MessageListVisibility: "show",
		Color:                 Block.Color,
	}
	// BlockDomain is the @SRC/Block/Domain folder for email domains that are automatically removed from your inbox.
	// Blocked emails are not analyzed by SRC.
	BlockDomain = gmail.Label{
		Name:                  Block.Name + "/Domain",
		MessageListVisibility: "show",
		Color:                 Block.Color,
	}
	// BlockGraveyard is the @SRC/Block/Graveyard folder for emails that have been removed from your inbox.
	// Blocked emails are not be analyzed by SRC.
	BlockGraveyard = gmail.Label{
		Name:                  Block.Name + "/🪦",
		MessageListVisibility: "show",
		Color:                 Block.Color,
	}
)
