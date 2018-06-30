package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type DataSource struct {
	Id     int64  `json:"id,omitempty"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	URL    string `json:"url"`
	Access string `json:"access"`

	Database string `json:"database,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`

	OrgId     int64 `json:"orgId,omitempty"`
	IsDefault bool  `json:"isDefault"`

	BasicAuth         bool   `json:"basicAuth"`
	BasicAuthUser     string `json:"basicAuthUser,omitempty"`
	BasicAuthPassword string `json:"basicAuthPassword,omitempty"`

	JSONData       JSONData       `json:"jsonData,omitempty"`
	SecureJSONData SecureJSONData `json:"secureJsonData,omitempty"`
}

// JSONData is a representation of the datasource `jsonData` property
type JSONData struct {
	AssumeRoleArn           string `json:"assumeRoleArn,omitempty"`
	AuthType                string `json:"authType,omitempty"`
	CustomMetricsNamespaces string `json:"customMetricsNamespaces,omitempty"`
	DefaultRegion           string `json:"defaultRegion,omitempty"`
}

// SecureJSONData is a representation of the datasource `secureJsonData` property
type SecureJSONData struct {
	AccessKey string `json:"accessKey,omitempty"`
	SecretKey string `json:"secretKey,omitempty"`
}

func (c *GrafanaClient) GetDataSources() (*[]DataSource, error) {
	dataSources := make([]DataSource, 0)
	req, err := c.newRequest("GET", "/api/datasources", nil)
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

	err = json.Unmarshal(data, &dataSources)
	return &dataSources, err
}
