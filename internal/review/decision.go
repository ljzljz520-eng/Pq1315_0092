package review

import (
	"aromaatelier/internal/model"
	"fmt"
)

type Decision struct {
	Approved bool
	Score    int
	Missing  []string
}

func (s *Service) Evaluate(recordID string, checklist model.ReviewChecklist) (Decision, error) {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return Decision{}, err
	}
	result := model.DecideReview(record, checklist)
	return Decision{Approved: result.Allowed, Score: result.Score, Missing: append([]string(nil), result.Reasons...)}, nil
}

func (s *Service) ApproveWithChecklist(recordID, actor, note string, checklist model.ReviewChecklist) (model.Record, error) {
	decision, err := s.Evaluate(recordID, checklist)
	if err != nil {
		return model.Record{}, err
	}
	if !decision.Approved {
		return model.Record{}, fmt.Errorf("review blocked: %v", decision.Missing)
	}
	return s.Approve(recordID, actor, note)
}

func (s *Service) QueueSummary() (map[string]int, error) {
	rows, err := s.store.ListRecords(model.SearchQuery{})
	if err != nil {
		return nil, err
	}
	result := map[string]int{"pending": 0, "approved": 0, "draft": 0, "archived": 0}
	for _, row := range rows {
		result[string(row.Status)]++
	}
	return result, nil
}

func (s *Service) CanArchive(recordID string) (bool, error) {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return false, err
	}
	return record.Status == model.StatusApproved, nil
}
