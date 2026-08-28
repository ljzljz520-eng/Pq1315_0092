package importer

import (
	"aromaatelier/internal/model"
	"fmt"
	"strings"
)

type Source struct {
	Name    string
	Headers []string
	Rows    []model.ImportRow
}

func ParseDelimited(input string, delimiter rune) (Source, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
		return Source{}, fmt.Errorf("header is required")
	}
	headers := strings.Split(lines[0], string(delimiter))
	if len(headers) < 4 {
		return Source{}, fmt.Errorf("at least four columns required")
	}
	rows := []model.ImportRow{}
	for index, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, string(delimiter))
		if len(parts) < 4 {
			return Source{}, fmt.Errorf("line %d has too few columns", index+2)
		}
		row := model.ImportRow{ExternalID: parts[0], Name: parts[1], Address: parts[2], Phone: parts[3]}
		if len(parts) > 4 {
			row.Email = parts[4]
		}
		if len(parts) > 5 {
			row.Tags = parts[5]
		}
		rows = append(rows, row)
	}
	return Source{Name: "delimited", Headers: headers, Rows: NormalizeRows(rows)}, nil
}

func (s Source) Valid() bool { return len(s.Headers) >= 4 && len(s.Rows) > 0 }

func (s Source) IDs() []string {
	result := make([]string, 0, len(s.Rows))
	for _, row := range s.Rows {
		result = append(result, row.ExternalID)
	}
	return result
}

func (s Source) Errors() []string {
	quality := Assess(s.Rows)
	result := []string{}
	for field, count := range quality.EmptyFields {
		result = append(result, fmt.Sprintf("%s empty in %d rows", field, count))
	}
	if len(quality.DuplicateIDs) > 0 {
		result = append(result, "duplicate ids: "+strings.Join(quality.DuplicateIDs, ","))
	}
	return result
}
