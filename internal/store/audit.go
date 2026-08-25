package store

import (
	"sort"

	"aromaatelier/internal/model"
	"go.etcd.io/bbolt"
)

func (s *DB) AppendAudit(recordID, action, actor, detail string, version int) (model.AuditEvent, error) {
	var event model.AuditEvent
	err := s.transaction(func(tx *bbolt.Tx) error {
		s.sequence++
		event = model.NewAudit(recordID, action, actor, detail, version, s.sequence)
		return putJSON(tx.Bucket([]byte("audit")), event.ID, event)
	})
	return event, err
}

func (s *DB) ListAudit(recordID string) ([]model.AuditEvent, error) {
	var events []model.AuditEvent
	err := s.view(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("audit")).ForEach(func(_, value []byte) error {
			if value == nil {
				return nil
			}
			var event model.AuditEvent
			if err := jsonUnmarshal(value, &event); err != nil {
				return err
			}
			if event.RecordID == recordID {
				events = append(events, event)
			}
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(events, func(i, j int) bool { return events[i].CreatedAt < events[j].CreatedAt })
	return events, nil
}

func (s *DB) AuditFor(recordID string) ([]model.AuditEvent, error) {
	var events []model.AuditEvent
	err := s.view(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("audit")).ForEach(func(key, value []byte) error {
			if value == nil {
				return nil
			}
			var event model.AuditEvent
			if err := jsonUnmarshal(value, &event); err != nil {
				return err
			}
			if event.RecordID == recordID {
				events = append(events, event)
			}
			return nil
		})
	})
	sort.Slice(events, func(i, j int) bool { return events[i].CreatedAt < events[j].CreatedAt })
	return events, err
}

func jsonUnmarshal(data []byte, target any) error {
	return unmarshal(data, target)
}
