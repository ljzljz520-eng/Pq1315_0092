package review

import (
	"aromaatelier/internal/model"
	"aromaatelier/internal/store"
	"fmt"
)

type Service struct{ store *store.DB }

func New(db *store.DB) *Service { return &Service{store: db} }

func (s *Service) Submit(recordID, actor string) (model.Record, error) {
	return s.transition(recordID, model.StatusPending, actor, "submitted for review")
}

func (s *Service) Approve(recordID, actor, note string) (model.Record, error) {
	if note == "" {
		return model.Record{}, fmt.Errorf("approval note is required")
	}
	return s.transition(recordID, model.StatusApproved, actor, note)
}

func (s *Service) Reject(recordID, actor, note string) (model.Record, error) {
	if note == "" {
		return model.Record{}, fmt.Errorf("rejection note is required")
	}
	return s.transition(recordID, model.StatusDraft, actor, note)
}

func (s *Service) transition(id string, target model.RecordStatus, actor, detail string) (model.Record, error) {
	record, err := s.store.GetRecord(id)
	if err != nil {
		return model.Record{}, err
	}
	if err := model.ValidateTransition(record.Status, target); err != nil {
		return model.Record{}, err
	}
	record.Status = target
	record.Version++
	record.UpdatedBy = actor
	if err := s.store.SaveRecord(record); err != nil {
		return model.Record{}, err
	}
	if _, err := s.store.AppendAudit(id, string(target), actor, detail, record.Version); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) Pending() ([]model.Record, error) {
	return s.store.ListRecords(model.SearchQuery{Status: model.StatusPending})
}

func (s *Service) Audit(recordID string) ([]model.AuditEvent, error) {
	return s.store.AuditFor(recordID)
}
