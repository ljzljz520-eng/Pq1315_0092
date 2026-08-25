package importer

import (
	"aromaatelier/internal/catalog"
	"aromaatelier/internal/model"
	"fmt"
	"strings"
)

type Result struct {
	Created  int
	Updated  int
	Rejected int
	Errors   []string
}

type Service struct{ catalog *catalog.Catalog }

func New(c *catalog.Catalog) *Service { return &Service{catalog: c} }

func (s *Service) Import(rows []model.ImportRow, actor string) Result {
	result := Result{}
	for index, row := range rows {
		if err := validateRow(row); err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: %v", index+1, err))
			continue
		}
		record := model.Record{ID: row.ExternalID, Name: row.Name, Address: row.Address, Phone: row.Phone, Email: row.Email, Tags: strings.Split(row.Tags, "|")}
		if _, err := s.catalog.Register(record, actor); err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: %v", index+1, err))
			continue
		}
		result.Created++
	}
	return result
}

func validateRow(row model.ImportRow) error {
	if strings.TrimSpace(row.ExternalID) == "" {
		return fmt.Errorf("external id is required")
	}
	if strings.TrimSpace(row.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(row.Address) == "" {
		return fmt.Errorf("address is required")
	}
	if strings.TrimSpace(row.Phone) == "" {
		return fmt.Errorf("phone is required")
	}
	return nil
}

func ParseTSV(input string) []model.ImportRow {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	rows := make([]model.ImportRow, 0, len(lines))
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) < 5 {
			continue
		}
		tags := ""
		if len(parts) > 5 {
			tags = parts[5]
		}
		rows = append(rows, model.ImportRow{ExternalID: parts[0], Name: parts[1], Address: parts[2], Phone: parts[3], Email: parts[4], Tags: tags})
	}
	return rows
}
