package kokoro

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/kamichidu/kokoro-go/types"
)

type Client struct {
	http.Client

	BaseUrl string

	tokenApplyer func(*http.Request) error
}

func NewClient(baseUrlStr string) *Client {
	return NewClientWithToken(baseUrlStr, "")
}

func NewClientWithToken(baseUrlStr string, token string) *Client {
	cli := new(Client)
	cli.BaseUrl = baseUrlStr
	rtc := rtchain{}
	rtc.Use(addContentTypeHeader)
	if token != "" {
		applyer := addXAccessTokenHeader(token)
		rtc.Use(applyer)
		cli.tokenApplyer = applyer
	}
	cli.Transport = rtc
	return cli
}

// Chainable http.RoundTripper
type rtchain []func(*http.Request) error

func (self *rtchain) Use(c func(*http.Request) error) {
	*self = append(*self, c)
}

func (self rtchain) RoundTrip(req *http.Request) (*http.Response, error) {
	for _, fn := range self {
		if err := fn(req); err != nil {
			return nil, err
		}
	}
	return http.DefaultTransport.RoundTrip(req)
}

func addXAccessTokenHeader(token string) func(*http.Request) error {
	return func(req *http.Request) error {
		if token == "" {
			return errors.New("no token to apply")
		}
		if req.Header.Get("X-Account-Token") != "" {
			return errors.New("X-Account-Token header is already exists")
		}
		req.Header.Set("X-Access-Token", token)
		return nil
	}
}

func addXAccountTokenHeader(email string, password string) func(*http.Request) error {
	return func(req *http.Request) error {
		if email == "" || password == "" {
			return errors.New("incomplete credentials")
		}
		if req.Header.Get("X-Access-Token") != "" {
			return errors.New("X-Access-Token header is already exists")
		}
		req.Header.Set("X-Account-Token", base64.StdEncoding.EncodeToString([]byte(email+":"+password)))
		return nil
	}
}

func addContentTypeHeader(req *http.Request) error {
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	return nil
}

func (self *Client) RegisterDevice(email string, password string, v *types.Device) (*types.Device, error) {
	body, err := makeJsonBody(v)
	if err != nil {
		return nil, fmt.Errorf("kokoro-go: %s", err)
	}

	u, err := makeRequestUrl(self.BaseUrl, "/api/v1/devices")
	if err != nil {
		return nil, fmt.Errorf("kokoro-go: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("kokoro-go: %s", err)
	}

	cli := new(http.Client)
	cli.Transport = rtchain{
		addContentTypeHeader,
		addXAccountTokenHeader(email, password),
	}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("kokoro-go: %s", err)
	}
	defer resp.Body.Close()

	if err := reportErrorResponse(resp); err != nil {
		return nil, err
	}

	retval := new(types.Device)
	return retval, json.NewDecoder(resp.Body).Decode(retval)
}

func (self *Client) doJson(method string, relativePath string, reqval interface{}, retval interface{}) error {
	u, err := makeRequestUrl(self.BaseUrl, relativePath)
	if err != nil {
		return fmt.Errorf("kokoro-go: %s", err)
	}
	body, err := makeJsonBody(reqval)
	if err != nil {
		return fmt.Errorf("kokoro-go: %s", err)
	}
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return fmt.Errorf("kokoro-go: %s", err)
	}

	resp, err := self.Do(req)
	if err != nil {
		return fmt.Errorf("kokoro-go: %s", err)
	}
	defer resp.Body.Close()

	if err := reportErrorResponse(resp); err != nil {
		return fmt.Errorf("kokoro-go: %s", err)
	}

	if err := json.NewDecoder(resp.Body).Decode(retval); err != nil {
		return fmt.Errorf("kokoro-go: %s", err)
	}
	return nil
}

func makeJsonBody(v interface{}) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}
	rw := new(bytes.Buffer)
	return rw, json.NewEncoder(rw).Encode(v)
}

func reportErrorResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	// TODO: report error detail
	return fmt.Errorf("%d - %s", resp.StatusCode, resp.Status)
}

func makeRequestUrl(baseUrlStr string, relativePath string) (*url.URL, error) {
	base, err := url.Parse(baseUrlStr)
	if err != nil {
		return nil, err
	}
	if base.Scheme == "" {
		base.Scheme = "https"
	}

	ref, err := url.ParseRequestURI(relativePath)
	if err != nil {
		return nil, err
	}
	return base.ResolveReference(ref), nil
}
