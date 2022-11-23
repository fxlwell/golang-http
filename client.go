package http

import (
	"bytes"
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

func Get(url string) *ClientRespone {
	return DefaultClient.Get(url)
}

func (c *Client) Get(url string) *ClientRespone {
	return c.doing(http.MethodGet, url, nil, nil)
}

func Post(url string, contentType string, headers http.Header, body io.Reader) *ClientRespone {
	return DefaultClient.Post(url, contentType, headers, body)
}

func (c *Client) Post(url string, contentType string, headers http.Header, body io.Reader) *ClientRespone {
	if headers == nil {
		headers = http.Header{}
	}
	headers.Set("Content-Type", contentType)
	return c.doing(http.MethodPost, url, headers, body)
}

func PostForm(httpUrl string, headers http.Header, data url.Values) *ClientRespone {
	return DefaultClient.PostForm(httpUrl, headers, data)
}

func (c *Client) PostForm(httpUrl string, headers http.Header, data url.Values) *ClientRespone {
	if data == nil {
		data = url.Values{}
	}
	return c.Post(httpUrl, "application/x-www-form-urlencoded", headers, strings.NewReader(data.Encode()))
}

func PostJsonBytes(httpUrl string, headers http.Header, jsonBytes []byte) *ClientRespone {
	return DefaultClient.PostJsonBytes(httpUrl, headers, jsonBytes)
}

func (c *Client) PostJsonBytes(httpUrl string, headers http.Header, jsonBytes []byte) *ClientRespone {
	return c.Post(httpUrl, "application/json", headers, bytes.NewReader(jsonBytes))
}

func PostJsonOBJ(httpUrl string, headers http.Header, jsonObj any) *ClientRespone {
	return DefaultClient.PostJsonOBJ(httpUrl, headers, jsonObj)
}

func (c *Client) PostJsonOBJ(httpUrl string, headers http.Header, jsonObj any) *ClientRespone {
	b, err := json.Marshal(jsonObj)
	if err != nil {
		return &ClientRespone{[]byte{}, nil, err}
	}
	return c.Post(httpUrl, "application/json", headers, bytes.NewReader(b))
}

func (c *Client) PostBodyString(httpUrl string, headers http.Header, data string) *ClientRespone {
	return c.Post(httpUrl, "application/x-www-form-urlencoded", headers, strings.NewReader(data))
}

func (c *Client) PostBodyBytes(httpUrl string, headers http.Header, data []byte) *ClientRespone {
	return c.Post(httpUrl, "application/x-www-form-urlencoded", headers, bytes.NewReader(data))
}
