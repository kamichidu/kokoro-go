package types

type Avatar struct {
	Size int `json:"size"`

	Url string `json:"url"`

	IsDefault bool `json:"is_default"`
}
