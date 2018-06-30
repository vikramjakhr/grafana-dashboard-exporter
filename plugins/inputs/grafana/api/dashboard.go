package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type DashboardMeta struct {
	IsStarred bool   `json:"isStarred"`
	URL       string `json:"url"`
}

type Dashboard struct {
	Meta  DashboardMeta          `json:"meta"`
	Model map[string]interface{} `json:"dashboard"`
}

func (c *GrafanaClient) GetDashboard(uri string) (*Dashboard, error) {
	path := fmt.Sprintf("/api/dashboards/%s", uri)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &Dashboard{}
	err = json.Unmarshal(data, &result)
	return result, err
}
