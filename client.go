package ehclient

import (
	"io"
	"net/http"
	"time"
)

type Endpoint string

var (
	EhEndpoint Endpoint = "e-hentai.org"
	ExEndpoint Endpoint = "exhentai.org"
)

type ClientOptions struct {
	Endpoint      Endpoint
	CookieExpires time.Duration
	UserAgent     string
}

var defaultOpts = &ClientOptions{
	Endpoint:      EhEndpoint,
	CookieExpires: time.Hour * 24 * 365,
}

type Client struct {
	httpc  *http.Client
	opts   *ClientOptions
	parser *Parser
}

func New(opt *ClientOptions) *Client {
	_opt := defaultOpts
	if opt != nil {
		if opt.Endpoint != "" {
			_opt.Endpoint = opt.Endpoint
		}
		if opt.CookieExpires != 0 {
			_opt.CookieExpires = opt.CookieExpires
		}
		if opt.UserAgent != "" {
			_opt.UserAgent = opt.UserAgent
		}
	}
	client := &Client{
		httpc:  &http.Client{},
		opts:   _opt,
		parser: NewParser(),
	}
	return client
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c.opts.UserAgent != "" {
		req.Header.Add("User-Agent", c.opts.UserAgent)
	}
	resp, err := c.httpc.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		bs, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, newHttpError(resp.StatusCode, resp.Status, "", err)
		}
		return nil, newHttpError(resp.StatusCode, resp.Status, string(bs), nil)
	}
	return resp, nil
}
