package catalog

import (
	"aromaatelier/internal/model"
	"sort"
	"strings"
)

type Facet struct {
	Value string
	Count int
}

type Insights struct {
	Total     int
	ByStatus  map[model.RecordStatus]int
	Tags      []Facet
	Addresses []Facet
}

func (c *Catalog) BuildInsights() (Insights, error) {
	rows, err := c.store.ListRecords(model.SearchQuery{})
	if err != nil {
		return Insights{}, err
	}
	result := Insights{ByStatus: map[model.RecordStatus]int{}}
	tagCounts := map[string]int{}
	addressCounts := map[string]int{}
	for _, row := range rows {
		result.Total++
		result.ByStatus[row.Status]++
		for _, tag := range row.Tags {
			tagCounts[tag]++
		}
		addressCounts[row.Address]++
	}
	result.Tags = facets(tagCounts)
	result.Addresses = facets(addressCounts)
	return result, nil
}

func facets(counts map[string]int) []Facet {
	result := make([]Facet, 0, len(counts))
	for value, count := range counts {
		result = append(result, Facet{Value: value, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Value < result[j].Value
		}
		return result[i].Count > result[j].Count
	})
	return result
}

func (c *Catalog) FindDuplicates() ([][]model.Record, error) {
	rows, err := c.store.ListRecords(model.SearchQuery{})
	if err != nil {
		return nil, err
	}
	groups := map[string][]model.Record{}
	for _, row := range rows {
		key := strings.ToLower(strings.TrimSpace(row.Name))
		groups[key] = append(groups[key], row)
	}
	duplicates := [][]model.Record{}
	for _, group := range groups {
		if len(group) > 1 {
			duplicates = append(duplicates, group)
		}
	}
	sort.Slice(duplicates, func(i, j int) bool { return duplicates[i][0].Name < duplicates[j][0].Name })
	return duplicates, nil
}

func MergeTags(existing, incoming []string) []string {
	combined := append(append([]string(nil), existing...), incoming...)
	return model.NormalizeTags(combined)
}

func (c *Catalog) ReplaceTags(id, actor string, tags []string) (model.Record, error) {
	return c.Edit(id, actor, func(record *model.Record) error { record.Tags = MergeTags(nil, tags); return nil })
}

func (c *Catalog) AddressHistory(id string) ([]string, error) {
	record, err := c.store.GetRecord(id)
	if err != nil {
		return nil, err
	}
	events, err := c.store.AuditFor(id)
	if err != nil {
		return nil, err
	}
	history := []string{}
	if record.PreviousAddr != "" {
		history = append(history, record.PreviousAddr)
	}
	for _, event := range events {
		if strings.Contains(event.Detail, "address") {
			history = append(history, event.Detail)
		}
	}
	return history, nil
}
