package store

import (
	"bytes"
	"sort"
	"strings"

	"aromaatelier/internal/model"
	"go.etcd.io/bbolt"
)

func (s *DB) SaveRecord(record model.Record) error {
	record = model.NormalizeRecord(record)
	if err := model.ValidateRecord(record); err != nil {
		return err
	}
	return s.transaction(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("records"))
		var existing model.Record
		if data := bucket.Get([]byte(record.ID)); data != nil {
			loaded, err := model.UnmarshalRecord(data)
			if err != nil {
				return err
			}
			existing = loaded
			if record.Version != existing.Version+1 {
				return ErrConflict
			}
		} else if record.Version != 1 {
			return ErrConflict
		}
		return putJSON(bucket, record.ID, record)
	})
}

func (s *DB) CreateRecord(record model.Record) error {
	record.Version = 1
	if record.Status == "" {
		record.Status = model.StatusDraft
	}
	return s.transaction(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("records"))
		if bucket.Get([]byte(record.ID)) != nil {
			return ErrConflict
		}
		record = model.NormalizeRecord(record)
		if err := model.ValidateRecord(record); err != nil {
			return err
		}
		return putJSON(bucket, record.ID, record)
	})
}

func (s *DB) GetRecord(id string) (model.Record, error) {
	var record model.Record
	err := s.view(func(tx *bbolt.Tx) error { return getJSON(tx.Bucket([]byte("records")), id, &record) })
	if err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *DB) ListRecords(query model.SearchQuery) ([]model.Record, error) {
	var records []model.Record
	err := s.view(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("records")).ForEach(func(_, value []byte) error {
			if value == nil {
				return nil
			}
			record, err := model.UnmarshalRecord(value)
			if err != nil {
				return err
			}
			if query.Status != "" && record.Status != query.Status {
				return nil
			}
			needle := strings.ToLower(strings.TrimSpace(query.Text))
			if needle != "" && !strings.Contains(strings.ToLower(record.Name+" "+record.Description+" "+record.Address), needle) {
				return nil
			}
			if query.Tag != "" && !contains(record.Tags, strings.ToLower(query.Tag)) {
				return nil
			}
			records = append(records, record)
			return nil
		})
	})
	sort.Slice(records, func(i, j int) bool { return bytes.Compare([]byte(records[i].ID), []byte(records[j].ID)) < 0 })
	if query.Limit > 0 && len(records) > query.Limit {
		records = records[:query.Limit]
	}
	return records, err
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func (s *DB) DeleteRecord(id string) error {
	return s.transaction(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("records"))
		if bucket.Get([]byte(id)) == nil {
			return ErrNotFound
		}
		return bucket.Delete([]byte(id))
	})
}
