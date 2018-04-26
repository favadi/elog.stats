package param

import (
	"github.com/txchuyen/elog.stats/elog/model"
	epb "github.com/txchuyen/elog.stats/pb/elog"
)

// Query ...
type Query struct {
	IPClient string
	IPServer string
	Tags     model.Tags
}

// Empty ...
func (q *Query) Empty() bool {
	return len(q.IPClient) == len(q.IPServer) &&
		len(q.IPClient) == 0 &&
		len(q.Tags) == 0
}

// FromPbQuery ...
func FromPbQuery(value *epb.Query) (*Query, error) {
	tags, err := model.FromPbTags(value.GetTags())
	if err != nil {
		return nil, err
	}
	return &Query{
		IPClient: value.GetIpClient(),
		IPServer: value.GetIpServer(),
		Tags:     tags,
	}, nil
}
