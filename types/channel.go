package types

type ChannelKind string

const (
	ChannelKindPublic        ChannelKind = "public_channel"
	ChannelKindPrivate                   = "private_channel"
	ChannelKindDirectMessage             = "direct_message"
)

type Channel struct {
	Id string `json:"id"`

	ChannelName string `json:"channel_name"`

	Kind ChannelKind `json:"kind"`

	Archived bool `json:"archived"`

	Description string `json:"description"`

	// Membership relation about this channel of login user
	Membership *Membership `json:"membership,omitempty"`

	// Membership relations about this channel
	Memberships []*Membership `json:"memberships,omitempty"`
}
