package flow031

import (
	"aromaatelier/internal/model"
	"sort"
	"strings"
)

type Dashboard struct {
	Total    int
	Draft    int
	Pending  int
	Approved int
	Archived int
	Tags     map[string]int
}

func BuildDashboard(records []model.Record) Dashboard {
	result := Dashboard{Tags: map[string]int{}}
	for _, record := range records {
		result.Total++
		switch record.Status {
		case model.StatusDraft:
			result.Draft++
		case model.StatusPending:
			result.Pending++
		case model.StatusApproved:
			result.Approved++
		case model.StatusArchived:
			result.Archived++
		}
		for _, tag := range record.Tags {
			result.Tags[tag]++
		}
	}
	return result
}

func DashboardRows(dashboard Dashboard) []string {
	return []string{"total=" + number(dashboard.Total), "draft=" + number(dashboard.Draft), "pending=" + number(dashboard.Pending), "approved=" + number(dashboard.Approved), "archived=" + number(dashboard.Archived)}
}

func TopTags(dashboard Dashboard, limit int) []string {
	type pair struct {
		name  string
		count int
	}
	pairs := []pair{}
	for name, count := range dashboard.Tags {
		pairs = append(pairs, pair{name, count})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count == pairs[j].count {
			return pairs[i].name < pairs[j].name
		}
		return pairs[i].count > pairs[j].count
	})
	if limit > 0 && len(pairs) > limit {
		pairs = pairs[:limit]
	}
	result := make([]string, 0, len(pairs))
	for _, item := range pairs {
		result = append(result, item.name)
	}
	return result
}

func (h *Handler) Dashboard(query model.SearchQuery) (Dashboard, error) {
	rows, err := h.Search(query)
	if err != nil {
		return Dashboard{}, err
	}
	return BuildDashboard(rows), nil
}

func (h *Handler) SearchNames(query model.SearchQuery) ([]string, error) {
	rows, err := h.Search(query)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.Name)
	}
	sort.Strings(result)
	return result, nil
}

func (h *Handler) TagSummary(query model.SearchQuery) (string, error) {
	rows, err := h.Search(query)
	if err != nil {
		return "", err
	}
	d := BuildDashboard(rows)
	tags := TopTags(d, 0)
	return strings.Join(tags, ", "), nil
}

func number(value int) string {
	if value == 0 {
		return "0"
	}
	negative := value < 0
	if negative {
		value = -value
	}
	result := ""
	for value > 0 {
		result = string(rune('0'+value%10)) + result
		value /= 10
	}
	if negative {
		return "-" + result
	}
	return result
}
