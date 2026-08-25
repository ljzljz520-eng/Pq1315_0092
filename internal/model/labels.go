package model

import (
	"sort"
	"strings"
)

type Label struct {
	Key   string
	Value string
}

func RecordLabels(record Record) []Label {
	labels := []Label{{Key: "status", Value: DisplayLabel(record.Status)}, {Key: "version", Value: formatInt(record.Version)}}
	if record.Email != "" {
		labels = append(labels, Label{Key: "email", Value: record.Email})
	}
	for _, tag := range SortedTags(record) {
		labels = append(labels, Label{Key: "tag", Value: tag})
	}
	return labels
}

func LabelMap(record Record) map[string][]string {
	result := map[string][]string{}
	for _, label := range RecordLabels(record) {
		result[label.Key] = append(result[label.Key], label.Value)
	}
	return result
}

func JoinLabels(labels []Label) string {
	parts := make([]string, 0, len(labels))
	for _, label := range labels {
		parts = append(parts, label.Key+"="+label.Value)
	}
	sort.Strings(parts)
	return strings.Join(parts, ";")
}

func formatInt(value int) string {
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
