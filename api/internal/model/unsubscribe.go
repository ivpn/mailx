package model

import "errors"

// One-click unsubscribe endpoints, each followed by the alias ID. They're called by the user's
// mail provider (RFC 8058) and authorised by a signed URL instead of a session.
const (
	OneClickDeactivatePath  = "/v1/deactivate-one-click/"
	OneClickDeletePath      = "/v1/delete-one-click/"
	OneClickBlockEmailPath  = "/v1/block-email-one-click/"
	OneClickBlockDomainPath = "/v1/block-domain-one-click/"
)

var (
	ErrOneClickAliasNotFound = errors.New("alias not found")
	ErrOneClickInvalidValue  = errors.New("invalid blocklist value")
	ErrOneClickAction        = errors.New("unable to complete one-click unsubscribe")
)

// OneClickRequest describes who triggered a one-click unsubscribe, for the user notification.
type OneClickRequest struct {
	IP        string
	UserAgent string
}
