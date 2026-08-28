package importer

import (
	"aromaatelier/internal/model"
	"sort"
	"strings"
)

type Quality struct {
	Total        int
	Valid        int
	Invalid      int
	DuplicateIDs []string
	EmptyFields  map[string]int
}

func Assess(rows []model.ImportRow) Quality {
	result := Quality{Total: len(rows), EmptyFields: map[string]int{}}
	seen := map[string]bool{}
	for _, row := range rows {
		valid := true
		if strings.TrimSpace(row.ExternalID) == "" {
			result.EmptyFields["id"]++
			valid = false
		}
		if strings.TrimSpace(row.Name) == "" {
			result.EmptyFields["name"]++
			valid = false
		}
		if strings.TrimSpace(row.Address) == "" {
			result.EmptyFields["address"]++
			valid = false
		}
		if strings.TrimSpace(row.Phone) == "" {
			result.EmptyFields["phone"]++
			valid = false
		}
		if seen[row.ExternalID] {
			result.DuplicateIDs = append(result.DuplicateIDs, row.ExternalID)
			valid = false
		}
		seen[row.ExternalID] = true
		if valid {
			result.Valid++
		} else {
			result.Invalid++
		}
	}
	sort.Strings(result.DuplicateIDs)
	return result
}

func NormalizeRows(rows []model.ImportRow) []model.ImportRow {
	result := make([]model.ImportRow, 0, len(rows))
	for _, row := range rows {
		row.ExternalID = strings.TrimSpace(row.ExternalID)
		row.Name = model.SanitizeText(row.Name)
		row.Address = model.SanitizeText(row.Address)
		row.Phone = model.SanitizeText(row.Phone)
		row.Email = strings.ToLower(strings.TrimSpace(row.Email))
		row.Tags = strings.ToLower(strings.TrimSpace(row.Tags))
		result = append(result, row)
	}
	return result
}

func ResultSummary(result Result) string {
	return strings.Join([]string{"created=" + itoa(result.Created), "updated=" + itoa(result.Updated), "rejected=" + itoa(result.Rejected)}, ", ")
}

func itoa(value int) string {
	digits := "0"
	if value < 0 {
		return "-" + itoa(-value)
	}
	if value == 0 {
		return digits
	}
	digits = ""
	for value > 0 {
		digits = string(rune('0'+value%10)) + digits
		value /= 10
	}
	return digits
}
