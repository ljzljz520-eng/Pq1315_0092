package model

import (
	"fmt"
	"sort"
	"strings"
)

type ReviewChecklist struct {
	IdentityVerified bool
	AddressVerified  bool
	SafetyNote       bool
	PhotoAttached    bool
	DescriptionReady bool
}

type ReviewDecision struct {
	Allowed bool
	Score   int
	Reasons []string
}

func (c ReviewChecklist) Score() int {
	score := 0
	if c.IdentityVerified {
		score += 2
	}
	if c.AddressVerified {
		score += 2
	}
	if c.SafetyNote {
		score += 2
	}
	if c.PhotoAttached {
		score++
	}
	if c.DescriptionReady {
		score++
	}
	return score
}

func (c ReviewChecklist) Missing() []string {
	missing := []string{}
	if !c.IdentityVerified {
		missing = append(missing, "identity")
	}
	if !c.AddressVerified {
		missing = append(missing, "address")
	}
	if !c.SafetyNote {
		missing = append(missing, "safety")
	}
	if !c.PhotoAttached {
		missing = append(missing, "photo")
	}
	if !c.DescriptionReady {
		missing = append(missing, "description")
	}
	return missing
}

func DecideReview(record Record, checklist ReviewChecklist) ReviewDecision {
	decision := ReviewDecision{Score: checklist.Score(), Reasons: []string{}}
	if record.Status != StatusPending {
		decision.Reasons = append(decision.Reasons, "record is not pending")
	}
	if len(record.Name) < 3 {
		decision.Reasons = append(decision.Reasons, "name is too short")
	}
	if len(record.Address) < 8 {
		decision.Reasons = append(decision.Reasons, "address is incomplete")
	}
	if missing := checklist.Missing(); len(missing) > 0 {
		decision.Reasons = append(decision.Reasons, "missing "+strings.Join(missing, ", "))
	}
	decision.Allowed = len(decision.Reasons) == 0 && decision.Score >= 7
	return decision
}

func SortedTags(record Record) []string {
	result := append([]string(nil), record.Tags...)
	sort.Strings(result)
	return result
}

func DisplayLabel(status RecordStatus) string {
	switch status {
	case StatusDraft:
		return "Draft"
	case StatusPending:
		return "Pending review"
	case StatusApproved:
		return "Approved"
	case StatusArchived:
		return "Archived"
	default:
		return fmt.Sprintf("Unknown (%s)", status)
	}
}

func IsPublishable(record Record) bool {
	return record.Status == StatusApproved && record.Address != "" && record.Name != ""
}
