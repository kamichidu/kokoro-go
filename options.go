package kokoro

import (
	"net/url"
)

type QueryOption func(q url.Values)

type QueryOptions []QueryOption

func (self QueryOptions) encode() string {
	if len(self) == 0 {
		return ""
	}
	q := url.Values{}
	for _, opt := range self {
		opt(q)
	}
	return "?" + q.Encode()
}
