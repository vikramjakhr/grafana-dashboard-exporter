package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Org struct {
	Id   int64
	Name string
}

func (c *GrafanaClient) GetCurrentOrg() (*Org, error) {
	org := &Org{}

	req, err := c.newRequest("GET", "/api/org/", nil)
	if err != nil {
		return org, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return org, err
	}
	if resp.StatusCode != 200 {
		return org, errors.New(resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return org, err
	}
	err = json.Unmarshal(data, &org)
	return org, err
}
