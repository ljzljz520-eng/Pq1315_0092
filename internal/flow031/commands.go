package flow031

import (
	"aromaatelier/internal/export"
	"aromaatelier/internal/model"
	"fmt"
	"sort"
)

type CommandResult struct {
	ID      string
	Status  string
	Message string
}

func (h *Handler) Publish(id, actor string, checklist model.ReviewChecklist) (CommandResult, error) {
	decision, err := h.review.Evaluate(id, checklist)
	if err != nil {
		return CommandResult{}, err
	}
	if !decision.Approved {
		return CommandResult{}, fmt.Errorf("cannot publish: %v", decision.Missing)
	}
	record, err := h.review.Approve(id, actor, "checklist passed")
	if err != nil {
		return CommandResult{}, err
	}
	return CommandResult{ID: record.ID, Status: string(record.Status), Message: "published"}, nil
}

func (h *Handler) ExportSummary(query model.SearchQuery) (string, error) {
	rows, err := h.Search(query)
	if err != nil {
		return "", err
	}
	return export.Summary(rows), nil
}

func (h *Handler) PublishedRecords(query model.SearchQuery) ([]model.Record, error) {
	rows, err := h.Search(query)
	if err != nil {
		return nil, err
	}
	return h.exporter.Published(rows), nil
}

func (h *Handler) AuditTimeline(id string) ([]string, error) {
	events, err := h.Audit(id)
	if err != nil {
		return nil, err
	}
	sort.Slice(events, func(i, j int) bool { return events[i].CreatedAt < events[j].CreatedAt })
	result := make([]string, 0, len(events))
	for _, event := range events {
		result = append(result, fmt.Sprintf("%d %s %s", event.Version, event.Action, event.Detail))
	}
	return result, nil
}

func (h *Handler) Markdown(id string) (string, error) { return h.exporter.Markdown(id) }

func (h *Handler) ValidateForArchive(id string) error {
	can, err := h.review.CanArchive(id)
	if err != nil {
		return err
	}
	if !can {
		return fmt.Errorf("record is not approved")
	}
	return nil
}
