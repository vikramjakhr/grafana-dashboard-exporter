package metric

import (
	"github.com/vikramjakhr/grafana-dashboard-exporter"
)

type metric struct {
	org    string
	mType   gde.ValueType
	title   string
	content []byte
}

func (m *metric) Org() string {
	return m.org
}

func (m *metric) Type() gde.ValueType {
	return m.mType
}

func (m *metric) Title() string {
	return m.title
}

func (m *metric) Content() string {
	return string(m.content)
}

func New(org string, mType gde.ValueType, title string, content []byte) *metric {
	return &metric{org: org, mType: mType, title: title, content: content}
}
