package model

import (
	"encoding/json"
	"strings"
)

type RecordStatus string

const (
	StatusDraft    RecordStatus = "draft"
	StatusPending  RecordStatus = "pending"
	StatusApproved RecordStatus = "approved"
	StatusArchived RecordStatus = "archived"
)

type Record struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Slug         string       `json:"slug"`
	Description  string       `json:"description"`
	Address      string       `json:"address"`
	PreviousAddr string       `json:"previous_address"`
	Phone        string       `json:"phone"`
	Email        string       `json:"email"`
	Status       RecordStatus `json:"status"`
	Version      int          `json:"version"`
	CreatedBy    string       `json:"created_by"`
	UpdatedBy    string       `json:"updated_by"`
	Tags         []string     `json:"tags"`
	Notes        []string     `json:"notes"`
}

type AuditEvent struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Action    string `json:"action"`
	Actor     string `json:"actor"`
	Detail    string `json:"detail"`
	Version   int    `json:"version"`
	CreatedAt int64  `json:"created_at"`
}

type Workflow struct {
	ID       string   `json:"id"`
	RecordID string   `json:"record_id"`
	Name     string   `json:"name"`
	State    string   `json:"state"`
	Steps    []string `json:"steps"`
	Current  int      `json:"current"`
	Owner    string   `json:"owner"`
	Revision int      `json:"revision"`
}

type Attachment struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Name      string `json:"name"`
	MediaType string `json:"media_type"`
	Content   string `json:"content"`
	Checksum  string `json:"checksum"`
}

type SearchQuery struct {
	Text   string
	Status RecordStatus
	Tag    string
	Limit  int
}

type ImportRow struct {
	ExternalID string
	Name       string
	Address    string
	Phone      string
	Email      string
	Tags       string
}

type ExportDocument struct {
	RecordID string       `json:"record_id"`
	Name     string       `json:"name"`
	Address  string       `json:"address"`
	Phone    string       `json:"phone"`
	Status   string       `json:"status"`
	Version  int          `json:"version"`
	Audit    []AuditEvent `json:"audit"`
}

func (r Record) Clone() Record {
	r.Tags = append([]string(nil), r.Tags...)
	r.Notes = append([]string(nil), r.Notes...)
	return r
}

func (r Record) Marshal() ([]byte, error) { return json.Marshal(r) }

func UnmarshalRecord(data []byte) (Record, error) {
	var r Record
	err := json.Unmarshal(data, &r)
	return r, err
}

func NormalizeTags(raw []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(raw))
	for _, item := range raw {
		value := strings.ToLower(strings.TrimSpace(item))
		if value != "" && !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func Statuses() []RecordStatus {
	return []RecordStatus{StatusDraft, StatusPending, StatusApproved, StatusArchived}
}
