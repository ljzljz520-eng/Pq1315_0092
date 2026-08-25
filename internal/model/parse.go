package model

import "strings"

func ParseStatus(value string) RecordStatus {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "draft":
		return StatusDraft
	case "pending", "review":
		return StatusPending
	case "approved", "published":
		return StatusApproved
	case "archived", "closed":
		return StatusArchived
	default:
		return RecordStatus(strings.ToLower(strings.TrimSpace(value)))
	}
}

func StatusRank(status RecordStatus) int {
	switch status {
	case StatusDraft:
		return 1
	case StatusPending:
		return 2
	case StatusApproved:
		return 3
	case StatusArchived:
		return 4
	default:
		return 0
	}
}

func IsTerminal(status RecordStatus) bool { return status == StatusArchived }

func NextStatus(status RecordStatus) RecordStatus {
	switch status {
	case StatusDraft:
		return StatusPending
	case StatusPending:
		return StatusApproved
	case StatusApproved:
		return StatusArchived
	default:
		return status
	}
}

func CompactDescription(record Record, width int) string {
	value := SanitizeText(record.Description)
	if width <= 0 || len(value) <= width {
		return value
	}
	if width < 4 {
		return value[:width]
	}
	return value[:width-3] + "..."
}

func HasTag(record Record, tag string) bool {
	wanted := strings.ToLower(strings.TrimSpace(tag))
	for _, candidate := range record.Tags {
		if candidate == wanted {
			return true
		}
	}
	return false
}
