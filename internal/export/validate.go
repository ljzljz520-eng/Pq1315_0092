package export

import (
	"aromaatelier/internal/model"
	"fmt"
	"strings"
)

type Validation struct {
	Valid  bool
	Issues []string
}

func ValidateDocument(doc model.ExportDocument) Validation {
	result := Validation{Valid: true, Issues: []string{}}
	if strings.TrimSpace(doc.RecordID) == "" {
		result.Issues = append(result.Issues, "record id")
		result.Valid = false
	}
	if strings.TrimSpace(doc.Name) == "" {
		result.Issues = append(result.Issues, "name")
		result.Valid = false
	}
	if strings.TrimSpace(doc.Address) == "" {
		result.Issues = append(result.Issues, "address")
		result.Valid = false
	}
	if doc.Version < 1 {
		result.Issues = append(result.Issues, "version")
		result.Valid = false
	}
	if doc.Status != string(model.StatusApproved) && doc.Status != string(model.StatusArchived) {
		result.Issues = append(result.Issues, "status")
		result.Valid = false
	}
	return result
}

func (s *Service) Validate(id string) (Validation, error) {
	doc, err := s.Document(id)
	if err != nil {
		return Validation{}, err
	}
	return ValidateDocument(doc), nil
}

func (s *Service) Filename(doc model.ExportDocument, extension string) string {
	safe := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(doc.Name), " ", "-"))
	if safe == "" {
		safe = doc.RecordID
	}
	if extension == "" {
		extension = "json"
	}
	return fmt.Sprintf("%s-v%d.%s", safe, doc.Version, extension)
}

func EscapeCSV(value string) string {
	if strings.ContainsAny(value, ",\"\n") {
		return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
	}
	return value
}

func CSVLine(record model.Record) string {
	return strings.Join([]string{EscapeCSV(record.ID), EscapeCSV(record.Name), EscapeCSV(record.Address), EscapeCSV(record.Phone), EscapeCSV(string(record.Status)), fmt.Sprintf("%d", record.Version)}, ",")
}
