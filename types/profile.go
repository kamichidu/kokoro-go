package types

type ProfileType string

const (
	ProfileTypeUser ProfileType = "user"
	ProfileTypeBot              = "bot"
)

type Profile struct {
	Id string `json:"id"`

	Type ProfileType `json:"type"`

	ScreenName string `json:"screen_name"`

	DisplayName string `json:"display_name"`

	Avatar string `json:"avatar"`

	Avatars []*Avatar `json:"avatars"`

	Archived bool `json:"archived"`

	InvitedChannelsCount int `json:"invited_channels_count"`
}
