package store

import (
	"aromaatelier/internal/model"
	"go.etcd.io/bbolt"
	"sort"
	"strings"
)

type Maintenance struct {
	Records     int
	AuditEvents int
	Workflows   int
	Attachments int
}

func (s *DB) Inspect() (Maintenance, error) {
	result := Maintenance{}
	err := s.view(func(tx *bbolt.Tx) error { return nil })
	if err != nil {
		return result, err
	}
	rows, err := s.ListRecords(model.SearchQuery{})
	if err != nil {
		return result, err
	}
	result.Records = len(rows)
	for _, row := range rows {
		events, eventErr := s.AuditFor(row.ID)
		if eventErr != nil {
			return result, eventErr
		}
		result.AuditEvents += len(events)
		attachments, attachmentErr := s.ListAttachments(row.ID)
		if attachmentErr != nil {
			return result, attachmentErr
		}
		result.Attachments += len(attachments)
		workflows, workflowErr := s.ListWorkflows(row.ID)
		if workflowErr != nil {
			return result, workflowErr
		}
		result.Workflows += len(workflows)
	}
	return result, nil
}

func (s *DB) IDs(status model.RecordStatus) ([]string, error) {
	rows, err := s.ListRecords(model.SearchQuery{Status: status})
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}
	sort.Strings(ids)
	return ids, nil
}

func (s *DB) HasRecord(id string) (bool, error) {
	_, err := s.GetRecord(id)
	if err == ErrNotFound {
		return false, nil
	}
	return err == nil, err
}

func (s *DB) FindByEmail(email string) ([]model.Record, error) {
	rows, err := s.ListRecords(model.SearchQuery{})
	if err != nil {
		return nil, err
	}
	needle := strings.ToLower(strings.TrimSpace(email))
	result := []model.Record{}
	for _, row := range rows {
		if strings.EqualFold(row.Email, needle) {
			result = append(result, row)
		}
	}
	return result, nil
}
