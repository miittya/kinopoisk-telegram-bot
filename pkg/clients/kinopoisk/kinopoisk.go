package kinopoisk

import (
	"encoding/json"
	"io"
	"kinopoisk-telegram-bot/pkg/config"
	"net/http"
	"net/url"
	"time"
)

const (
	pagesNum  = "1"
	moviesNum = "5"
)

type Client struct {
	client *http.Client
	cfg    *config.Config
}

func NewClient(cfg *config.Config) (*Client, error) {
	return &Client{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		cfg: cfg,
	}, nil
}

func (c *Client) Get(title string) (res Response, err error) {
	q := url.Values{}
	q.Add("page", pagesNum)
	q.Add("limit", moviesNum)
	q.Add("query", title)
	req, err := c.generateRequest(c.cfg.EndPointMovieSearch, q)
	if err != nil {
		return Response{}, err
	}
	result, err := c.doHTTP(req)
	if err != nil {
		return Response{}, err
	}
	return result, nil
}

func (c *Client) generateRequest(endpoint string, query url.Values) (*http.Request, error) {
	u, err := url.Parse(c.cfg.APIHost + endpoint)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("X-API-KEY", c.cfg.KinopoiskToken)
	req.URL.RawQuery = query.Encode()
	return req, nil
}

func (c *Client) doHTTP(req *http.Request) (Response, error) {
	resp, err := http.DefaultClient.Do(req)
	defer func() { _ = resp.Body.Close() }()
	if err != nil {
		return Response{}, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return Response{}, err
	}
	return result, nil
}
