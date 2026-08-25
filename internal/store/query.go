package store

import (
	"aromaatelier/internal/model"
	"go.etcd.io/bbolt"
	"strings"
)

func (s *DB) CountByStatus() (map[model.RecordStatus]int, error) {
	counts := map[model.RecordStatus]int{}
	err := s.view(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("records")).ForEach(func(_, value []byte) error {
			if value == nil {
				return nil
			}
			var record model.Record
			if err := unmarshal(value, &record); err != nil {
				return err
			}
			counts[record.Status]++
			return nil
		})
	})
	return counts, err
}

func (s *DB) SearchByPrefix(prefix string) ([]model.Record, error) {
	needle := strings.ToLower(strings.TrimSpace(prefix))
	result := []model.Record{}
	rows, err := s.ListRecords(model.SearchQuery{})
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		if strings.HasPrefix(strings.ToLower(row.Name), needle) || strings.HasPrefix(strings.ToLower(row.ID), needle) {
			result = append(result, row)
		}
	}
	return result, nil
}

func (s *DB) ReplaceRecordAddress(id, address string, expectedVersion int) (model.Record, error) {
	record, err := s.GetRecord(id)
	if err != nil {
		return model.Record{}, err
	}
	if record.Version != expectedVersion {
		return model.Record{}, ErrConflict
	}
	if strings.TrimSpace(address) == "" {
		return model.Record{}, modelError("address is required")
	}
	record.PreviousAddr = record.Address
	record.Address = model.SanitizeText(address)
	record.Version++
	if err := s.SaveRecord(record); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

type modelError string

func (e modelError) Error() string { return string(e) }

func (s *DB) Snapshot() ([]model.Record, error) { return s.ListRecords(model.SearchQuery{}) }
