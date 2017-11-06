package types

type Authority string

const (
	AuthorityAdministrator Authority = "administrator"
	AuthorityMaintainer              = "maintainer"
	AuthorityMember                  = "member"
	AuthorityInvited                 = "invited"
)

type Membership struct {
	Id string `json:"id"`

	Channel *Channel `json:"channel,omitempty"`

	Authority Authority `json:"authority"`

	DisableNotification bool `json:"disable_notification"`

	UnreadCount int `json:"unread_count"`

	Profile *Profile `json:"profile"`
}
