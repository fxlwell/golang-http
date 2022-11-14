package http

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ClientConf struct {
	ConnTimeout           time.Duration
	KeepAlive             time.Duration
	IdleConnTimeout       time.Duration
	TlsHandshakeTimeout   time.Duration
	ExpectContinueTimeout time.Duration
	RequestTimeout        time.Duration
	InsecureSkipVerify    bool
	MaxIdleConns          int
}

var DefaultClientConfig = &ClientConf{
	ConnTimeout:           1 * time.Second,
	KeepAlive:             30 * time.Second,
	MaxIdleConns:          8,
	IdleConnTimeout:       90 * time.Second,
	TlsHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	InsecureSkipVerify:    false,
	RequestTimeout:        5 * time.Second,
}

var (
	DefaultClient *Client
	ErrNot200     = errors.New("not 200 ok")
)

func init() {
	DefaultClient = NewClient(DefaultClientConfig)
}

type Client struct {
	*http.Client
	conf *ClientConf
}

type ClientRespone struct {
	body []byte
	resp *http.Response
	err  error
}

func (cr *ClientRespone) Bytes() ([]byte, *http.Response, error) {
	return cr.body, cr.resp, cr.err
}

func (cr *ClientRespone) String() (string, *http.Response, error) {
	return string(cr.body), cr.resp, cr.err
}

func (cr *ClientRespone) JsonOBJ(obj any) (string, *http.Response, error) {
	if err := json.Unmarshal(cr.body, obj); err != nil {
		if cr.err != nil {
			return string(cr.body), cr.resp, fmt.Errorf("%w|%s", cr.err, err.Error())
		} else {
			return string(cr.body), cr.resp, err
		}
	}
	return string(cr.body), cr.resp, cr.err
}

func NewClient(conf *ClientConf) *Client {
	if conf == nil {
		conf = DefaultClientConfig
	}

	return &Client{
		&http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   conf.ConnTimeout,
					KeepAlive: conf.KeepAlive,
					DualStack: false,
				}).DialContext,
				MaxIdleConns:          conf.MaxIdleConns,
				IdleConnTimeout:       conf.IdleConnTimeout,
				TLSHandshakeTimeout:   conf.TlsHandshakeTimeout,
				ExpectContinueTimeout: conf.ExpectContinueTimeout,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: conf.InsecureSkipVerify,
				},
			},
			Timeout: conf.RequestTimeout,
		},
		conf,
	}
}

func (c *Client) doing(method, url string, headers http.Header, body io.Reader) *ClientRespone {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return &ClientRespone{[]byte{}, nil, err}
	}

	if headers != nil {
		for name, values := range headers {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}
	}

	resp, err := c.Do(req)
	if err != nil {
		return &ClientRespone{[]byte{}, resp, err}
	}
	defer resp.Body.Close()

	var b []byte

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return &ClientRespone{[]byte{}, resp, err}
	}

	if resp.StatusCode != http.StatusOK {
		return &ClientRespone{b, resp, ErrNot200}
	}

	return &ClientRespone{b, resp, nil}
}

func (c *Client) Get(url string) *ClientRespone {
	return c.doing(http.MethodGet, url, nil, nil)
}

func (c *Client) Post(url string, data url.Values) *ClientRespone {
	headers := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.doing(http.MethodPost, url, headers, strings.NewReader(data.Encode()))
}

func (c *Client) PostWithHeaders(url string, data url.Values, headers http.Header) *ClientRespone {
	if len(headers.Get("Content-Type")) <= 0 {
		headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return c.doing(http.MethodPost, url, headers, strings.NewReader(data.Encode()))
}
