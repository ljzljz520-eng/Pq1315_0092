package catalog

import (
	"aromaatelier/internal/model"
	"fmt"
	"sort"
	"strings"
)

type ReportLine struct {
	RecordID  string
	Name      string
	Status    string
	Address   string
	Attention string
}

func (c *Catalog) ReviewReport() ([]ReportLine, error) {
	rows, err := c.store.ListRecords(model.SearchQuery{})
	if err != nil {
		return nil, err
	}
	result := make([]ReportLine, 0, len(rows))
	for _, row := range rows {
		attention := "ready"
		switch row.Status {
		case model.StatusDraft:
			attention = "needs submission"
		case model.StatusPending:
			attention = "awaiting review"
		case model.StatusApproved:
			attention = "ready to publish"
		case model.StatusArchived:
			attention = "closed"
		}
		result = append(result, ReportLine{RecordID: row.ID, Name: row.Name, Status: model.DisplayLabel(row.Status), Address: row.Address, Attention: attention})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result, nil
}

func RenderReport(lines []ReportLine) string {
	parts := []string{"record | name | status | address | attention"}
	for _, line := range lines {
		parts = append(parts, fmt.Sprintf("%s | %s | %s | %s | %s", line.RecordID, line.Name, line.Status, line.Address, line.Attention))
	}
	return strings.Join(parts, "\n")
}

func ValidateForExport(record model.Record) []string {
	issues := []string{}
	if strings.TrimSpace(record.Name) == "" {
		issues = append(issues, "name")
	}
	if strings.TrimSpace(record.Address) == "" {
		issues = append(issues, "address")
	}
	if strings.TrimSpace(record.Phone) == "" {
		issues = append(issues, "phone")
	}
	if record.Status != model.StatusApproved {
		issues = append(issues, "approval")
	}
	return issues
}

func FilterByTag(records []model.Record, tag string) []model.Record {
	needle := strings.ToLower(strings.TrimSpace(tag))
	result := []model.Record{}
	for _, record := range records {
		for _, candidate := range record.Tags {
			if candidate == needle {
				result = append(result, record.Clone())
				break
			}
		}
	}
	return result
}
