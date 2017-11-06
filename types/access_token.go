package types

type AccessTokenKind string

const (
	AccessTokenKindUser      AccessTokenKind = "user"
	AccessTokenKindDevice                    = "device"
	AccessTokenKindEssential                 = "essential"
)

type AccessToken struct {
	Id string `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Token string `json:"token,omitempty"`

	Kind AccessTokenKind `json:"kind,omitempty"`
}
