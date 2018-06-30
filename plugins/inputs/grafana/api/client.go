package api

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type GrafanaClient struct {
	key     string
	baseURL url.URL
	*http.Client
}

//New creates a new grafana client
//auth can be in user:pass format, or it can be an api key
func NewGrafanaClient(auth, baseURL string) (*GrafanaClient, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	key := ""
	if strings.Contains(auth, ":") {
		split := strings.Split(auth, ":")
		u.User = url.UserPassword(split[0], split[1])
	} else {
		key = fmt.Sprintf("Bearer %s", auth)
	}
	return &GrafanaClient{
		key,
		*u,
		&http.Client{},
	}, nil
}

func (c *GrafanaClient) newRequest(method, requestPath string, body io.Reader) (*http.Request, error) {
	gURL := c.baseURL
	gURL.Path = path.Join(gURL.Path, requestPath)
	req, err := http.NewRequest(method, gURL.String(), body)
	if err != nil {
		return req, err
	}
	if c.key != "" {
		req.Header.Add("Authorization", c.key)
	}

	if body == nil {
		log.Printf("D! Request to %s with empty body data", gURL.String())
	} else {
		log.Printf("D! Request to %s with body data %s", gURL.String(), body.(*bytes.Buffer).String())
	}

	req.Header.Add("Content-Type", "application/json")
	return req, err
}
