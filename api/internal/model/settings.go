package model

type Settings struct {
	BaseModel
	UserID                   string                   `json:"-"`
	Domain                   string                   `json:"domain"`
	Recipient                string                   `json:"recipient"`
	FromName                 string                   `json:"from_name"`
	AliasFormat              string                   `json:"alias_format"`
	LogIssues                bool                     `json:"log_issues"`
	RemoveHeader             bool                     `json:"remove_header"`
	ListUnsubscribeBehaviour ListUnsubscribeBehaviour `json:"list_unsubscribe_behaviour" gorm:"type:smallint;not null;default:0"`
}

// ListUnsubscribeBehaviour controls the List-Unsubscribe/List-Unsubscribe-Post headers set on
// forwarded mail, i.e. what a mail client's "Unsubscribe" button does.
type ListUnsubscribeBehaviour int

const (
	// ListUnsubscribeOriginalWithFallback keeps the sender's header (mailto: rewritten to go
	// out through the alias), falling back to one-click deactivate when there is none.
	ListUnsubscribeOriginalWithFallback ListUnsubscribeBehaviour = iota
	ListUnsubscribeDeactivate
	ListUnsubscribeDelete
	ListUnsubscribeBlockEmail
	ListUnsubscribeBlockDomain
)

func (b ListUnsubscribeBehaviour) Valid() bool {
	return b >= ListUnsubscribeOriginalWithFallback && b <= ListUnsubscribeBlockDomain
}

// ListUnsubscribeHeaders are the List-Unsubscribe header values set on a forwarded message;
// an empty Value means neither header is set.
type ListUnsubscribeHeaders struct {
	Value string
	Post  string
}
