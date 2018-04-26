package store

import (
	"github.com/go-pg/pg"
	"github.com/txchuyen/elog.stats/elog/model"
	"github.com/txchuyen/elog.stats/elog/param"
)

// Eventer ...
type Eventer interface {
	Create(*model.Event) error
	List(*param.Query) ([]*model.Event, error)
}

type eventStore struct {
	db *pg.DB
}

// NewEvent create new event store
func NewEvent(db *pg.DB) Eventer {
	return &eventStore{
		db: db,
	}
}

// Create ...
func (s *eventStore) Create(e *model.Event) error {
	return s.db.Insert(e)
}

// List ...
func (s *eventStore) List(q *param.Query) (events []*model.Event, err error) {
	query := s.db.Model(&events)
	if q == nil || q.Empty() {
		err = query.Select()
		return
	}
	if address := q.IPClient; len(address) > 0 {
		query = query.Where("ip_client = ?", address)
	}
	if address := q.IPServer; len(address) > 0 {
		query = query.Where("ip_server = ?", address)
	}
	if tags := q.Tags; len(tags) > 0 {
		query = query.Where("tags @> ?", tags)
	}
	err = query.Select()
	return events, err
}
