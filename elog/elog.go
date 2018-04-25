package elog

import (
	"context"

	"elog.stats/elog/model"
	"elog.stats/elog/param"
	"elog.stats/elog/store"

	epb "elog.stats/pb/elog"
	gpb_empty "github.com/golang/protobuf/ptypes/empty"
)

type elogServer struct {
	Store store.Eventer
}

// NewServer create new instance of elog
func NewServer(store store.Eventer) epb.ElogServer {
	return &elogServer{
		Store: store,
	}
}

func (e *elogServer) List(request *epb.Query, stream epb.Elog_ListServer) error {
	query, err := param.FromPbQuery(request)
	if err != nil {
		return err
	}
	events, err := e.Store.List(query)
	if err != nil {
		return err
	}
	for _, e := range events {
		pbEvent, err := model.ToPbEvent(e)
		if err != nil {
			return err
		}
		if err = stream.Send(pbEvent); err != nil {
			return err
		}
	}
	return nil
}

func (e *elogServer) Create(ctx context.Context, request *epb.Event) (*gpb_empty.Empty, error) {
	event, err := model.FromPbEvent(request)
	if err != nil {
		return nil, err
	}
	err = e.Store.Create(event)
	if err != nil {
		return nil, err
	}
	return &gpb_empty.Empty{}, nil
}
