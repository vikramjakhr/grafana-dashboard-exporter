package metric

import (
	"github.com/vikramjakhr/grafana-dashboard-exporter"
)

type metric struct {
	name    string
	mType   gde.ValueType
	content []byte
}

func (m *metric) Name() string {
	return m.name
}

func (m *metric) Type() gde.ValueType {
	return m.mType
}

func (m *metric) Content() string {
	return string(m.content)
}

func New(name string, mType gde.ValueType, content []byte) *metric {
	return &metric{name: name, mType: mType, content: content}
}
