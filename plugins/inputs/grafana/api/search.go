package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type SearchResp struct {
	Id    int64
	Title string
	Uri   string
	Type  string
}

var (
	SearchTypeDashDB = "dash-db"
)

func (c *GrafanaClient) Search(sType, query string) (*[]SearchResp, error) {
	result := make([]SearchResp, 0)

	req, err := c.newRequest("GET", "/api/search", nil)
	if err != nil {
		return &result, err
	}

	q := req.URL.Query()
	q.Add("type", sType)
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()

	resp, err := c.Do(req)
	if err != nil {
		return &result, err
	}
	if resp.StatusCode != 200 {
		return &result, errors.New(resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &result, err
	}
	err = json.Unmarshal(data, &result)
	return &result, err
}
