package model

import (
	"errors"

	epb "elog.stats/pb/elog"
)

// Event ...
type Event struct {
	IPClient string
	IPServer string
	Tags     Tags
	Message  string
}

// FromPbEvent ...
func FromPbEvent(value *epb.Event) (*Event, error) {
	t, err := FromPbTags(value.GetTags())
	if err != nil {
		return nil, err
	}
	msg := value.GetMessage()
	if len(msg) == 0 {
		return nil, errors.New(`event: "message" is required`)
	}
	return &Event{
		IPClient: value.GetIpClient(),
		IPServer: value.GetIpServer(),
		Tags:     t,
		Message:  value.GetMessage(),
	}, nil
}

// ToPbEvent ...
func ToPbEvent(value *Event) (*epb.Event, error) {
	tags, err := ToPbTags(value.Tags)
	if err != nil {
		return nil, err
	}
	return &epb.Event{
		IpClient: value.IPClient,
		IpServer: value.IPServer,
		Tags:     tags,
		Message:  value.Message,
	}, nil
}
