package model

import (
	"fmt"
	"net/mail"
	"strings"
)

func ValidateRecord(r Record) error {
	if strings.TrimSpace(r.ID) == "" {
		return fmt.Errorf("id is required")
	}
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if len(r.Name) > 120 {
		return fmt.Errorf("name exceeds 120 characters")
	}
	if strings.TrimSpace(r.Address) == "" {
		return fmt.Errorf("address is required")
	}
	if len(r.Address) > 240 {
		return fmt.Errorf("address exceeds 240 characters")
	}
	if strings.TrimSpace(r.Phone) == "" {
		return fmt.Errorf("phone is required")
	}
	if strings.TrimSpace(r.Email) != "" {
		if _, err := mail.ParseAddress(r.Email); err != nil {
			return fmt.Errorf("invalid email: %w", err)
		}
	}
	if r.Status == "" {
		return fmt.Errorf("status is required")
	}
	if !IsKnownStatus(r.Status) {
		return fmt.Errorf("unknown status %q", r.Status)
	}
	if r.Version < 1 {
		return fmt.Errorf("version must be positive")
	}
	if len(r.Tags) > 20 {
		return fmt.Errorf("too many tags")
	}
	return nil
}

func IsKnownStatus(status RecordStatus) bool {
	for _, candidate := range Statuses() {
		if status == candidate {
			return true
		}
	}
	return false
}

func ValidateTransition(from, to RecordStatus) error {
	if from == to {
		return fmt.Errorf("status is unchanged")
	}
	switch from {
	case StatusDraft:
		if to != StatusPending {
			return fmt.Errorf("draft can only move to pending")
		}
	case StatusPending:
		if to != StatusApproved && to != StatusDraft {
			return fmt.Errorf("pending can move to approved or draft")
		}
	case StatusApproved:
		if to != StatusArchived && to != StatusPending {
			return fmt.Errorf("approved can move to archived or pending")
		}
	case StatusArchived:
		return fmt.Errorf("archived record cannot transition")
	default:
		return fmt.Errorf("unknown source status")
	}
	return nil
}

func SanitizeText(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func NormalizeRecord(r Record) Record {
	r.Name = SanitizeText(r.Name)
	r.Description = SanitizeText(r.Description)
	r.Address = SanitizeText(r.Address)
	r.Phone = SanitizeText(r.Phone)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.Tags = NormalizeTags(r.Tags)
	return r
}
