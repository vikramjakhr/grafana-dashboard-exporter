package metric

import (
	"github.com/vikramjakhr/grafana-dashboard-exporter"
)

type metric struct {
	dir    string
	mType   gde.ValueType
	title   string
	content []byte
}

func (m *metric) Dir() string {
	return m.dir
}

func (m *metric) Type() gde.ValueType {
	return m.mType
}

func (m *metric) Title() string {
	return m.title
}

func (m *metric) Content() []byte {
	return m.content
}

func New(dir string, mType gde.ValueType, title string, content []byte) *metric {
	return &metric{dir: dir, mType: mType, title: title, content: content}
}
