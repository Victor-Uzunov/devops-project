package constants

import "time"

const (
	ContentTypeJSON     = "application/json"
	AuthorizationHeader = "Authorization"
	TokenCtxKey         = "token"
	AdminOrganization   = "Admin-Role"
	WriterOrganization  = "Writer-Role"
	ReaderOrganization  = "Reader-Role"
	StatusAccepted      = "accepted"
	StatusOwner         = "owner"
	StatusPending       = "pending"
	DefaultDateTime     = "01-01-0001"
	CookieAge           = 31536000
	DateFormat          = time.RFC3339
)
