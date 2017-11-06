package kokoro

import (
	"net/http"
	"net/url"
	"path"

	"github.com/kamichidu/kokoro-go/types"
)

func (self *Client) CreateChannel(v *types.Channel) (*types.Channel, error) {
	val := &types.Channel{}
	return val, self.doJson(http.MethodPost, "/api/v1/channels", v, val)
}

var (
	Archived    QueryOption = func(q url.Values) { q.Set("archived", "true") }
	NotArchived             = func(q url.Values) { q.Set("archived", "false") }
)

func (self *Client) ListChannels(opts ...QueryOption) ([]*types.Channel, error) {
	val := []*types.Channel{}
	return val, self.doJson(http.MethodGet, "/api/v1/channels"+QueryOptions(opts).encode(), nil, &val)
}

func (self *Client) GetChannel(channelId string) (*types.Channel, error) {
	val := &types.Channel{}
	return val, self.doJson(http.MethodGet, path.Join("/api/v1/channels", channelId), nil, val)
}
