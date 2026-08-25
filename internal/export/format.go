package export

import (
	"aromaatelier/internal/model"
	"fmt"
	"sort"
	"strings"
)

func Markdown(doc model.ExportDocument) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", doc.Name)
	fmt.Fprintf(&b, "- Status: %s\n- Address: %s\n- Phone: %s\n- Version: %d\n\n", model.DisplayLabel(model.RecordStatus(doc.Status)), doc.Address, doc.Phone, doc.Version)
	if len(doc.Audit) > 0 {
		b.WriteString("## History\n\n")
		for _, event := range doc.Audit {
			fmt.Fprintf(&b, "- %s: %s (%s)\n", event.Action, event.Detail, event.Actor)
		}
	}
	return b.String()
}

func Summary(records []model.Record) string {
	copyRows := append([]model.Record(nil), records...)
	sort.Slice(copyRows, func(i, j int) bool { return copyRows[i].Name < copyRows[j].Name })
	lines := []string{fmt.Sprintf("records: %d", len(copyRows))}
	for _, record := range copyRows {
		lines = append(lines, fmt.Sprintf("%s | %s | %s", record.Name, model.DisplayLabel(record.Status), record.Address))
	}
	return strings.Join(lines, "\n")
}

func (s *Service) Markdown(id string) (string, error) {
	doc, err := s.Document(id)
	if err != nil {
		return "", err
	}
	return Markdown(doc), nil
}

func (s *Service) Published(records []model.Record) []model.Record {
	result := []model.Record{}
	for _, record := range records {
		if model.IsPublishable(record) {
			result = append(result, record.Clone())
		}
	}
	return result
}
