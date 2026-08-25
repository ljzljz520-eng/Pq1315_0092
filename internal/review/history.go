package review

import (
	"aromaatelier/internal/model"
	"sort"
	"strings"
)

type HistorySummary struct {
	RecordID   string
	FirstActor string
	LastActor  string
	Actions    []string
	Versions   int
}

func SummarizeHistory(events []model.AuditEvent) HistorySummary {
	result := HistorySummary{Actions: []string{}}
	if len(events) == 0 {
		return result
	}
	sort.Slice(events, func(i, j int) bool { return events[i].CreatedAt < events[j].CreatedAt })
	result.RecordID = events[0].RecordID
	result.FirstActor = events[0].Actor
	result.LastActor = events[len(events)-1].Actor
	seen := map[string]bool{}
	for _, event := range events {
		if !seen[event.Action] {
			result.Actions = append(result.Actions, event.Action)
			seen[event.Action] = true
		}
		if event.Version > result.Versions {
			result.Versions = event.Version
		}
	}
	return result
}

func (s *Service) HistorySummary(recordID string) (HistorySummary, error) {
	events, err := s.Audit(recordID)
	if err != nil {
		return HistorySummary{}, err
	}
	return SummarizeHistory(events), nil
}

func (s *Service) HasAction(recordID, action string) (bool, error) {
	events, err := s.Audit(recordID)
	if err != nil {
		return false, err
	}
	for _, event := range events {
		if strings.EqualFold(event.Action, action) {
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) LatestVersion(recordID string) (int, error) {
	summary, err := s.HistorySummary(recordID)
	if err != nil {
		return 0, err
	}
	return summary.Versions, nil
}

func (s *Service) PendingIDs() ([]string, error) {
	rows, err := s.Pending()
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.ID)
	}
	sort.Strings(result)
	return result, nil
}
