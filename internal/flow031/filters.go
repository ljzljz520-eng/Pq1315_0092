package flow031

import (
	"aromaatelier/internal/model"
	"sort"
	"strings"
)

type Filter struct {
	Status        model.RecordStatus
	Tag           string
	Prefix        string
	PublishedOnly bool
}

func ApplyFilter(records []model.Record, filter Filter) []model.Record {
	result := []model.Record{}
	for _, record := range records {
		if filter.Status != "" && record.Status != filter.Status {
			continue
		}
		if filter.Tag != "" && !model.HasTag(record, filter.Tag) {
			continue
		}
		if filter.Prefix != "" && !strings.HasPrefix(strings.ToLower(record.Name), strings.ToLower(filter.Prefix)) {
			continue
		}
		if filter.PublishedOnly && !model.IsPublishable(record) {
			continue
		}
		result = append(result, record.Clone())
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}

func (h *Handler) Filter(query model.SearchQuery, filter Filter) ([]model.Record, error) {
	rows, err := h.Search(query)
	if err != nil {
		return nil, err
	}
	return ApplyFilter(rows, filter), nil
}

func (h *Handler) Count(query model.SearchQuery, filter Filter) (int, error) {
	rows, err := h.Filter(query, filter)
	if err != nil {
		return 0, err
	}
	return len(rows), nil
}

func MatchAll(record model.Record, filters []Filter) bool {
	for _, filter := range filters {
		if len(ApplyFilter([]model.Record{record}, filter)) == 0 {
			return false
		}
	}
	return true
}

func FilterSummary(filter Filter) string {
	parts := []string{}
	if filter.Status != "" {
		parts = append(parts, "status="+string(filter.Status))
	}
	if filter.Tag != "" {
		parts = append(parts, "tag="+filter.Tag)
	}
	if filter.Prefix != "" {
		parts = append(parts, "prefix="+filter.Prefix)
	}
	if filter.PublishedOnly {
		parts = append(parts, "published=true")
	}
	return strings.Join(parts, ",")
}

func SelectLatest(records []model.Record) (model.Record, bool) {
	if len(records) == 0 {
		return model.Record{}, false
	}
	latest := records[0]
	for _, record := range records[1:] {
		if record.Version > latest.Version {
			latest = record
		}
	}
	return latest, true
}

func StatusCounts(records []model.Record) map[model.RecordStatus]int {
	counts := map[model.RecordStatus]int{}
	for _, record := range records {
		counts[record.Status]++
	}
	return counts
}

func IDs(records []model.Record) []string {
	result := make([]string, 0, len(records))
	for _, record := range records {
		result = append(result, record.ID)
	}
	sort.Strings(result)
	return result
}

func ContainsID(records []model.Record, id string) bool {
	for _, record := range records {
		if record.ID == id {
			return true
		}
	}
	return false
}

func Unique(records []model.Record) []model.Record {
	seen := map[string]bool{}
	result := []model.Record{}
	for _, record := range records {
		if !seen[record.ID] {
			seen[record.ID] = true
			result = append(result, record.Clone())
		}
	}
	return result
}
