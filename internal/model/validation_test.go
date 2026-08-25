package model

import "testing"

func TestNormalizeRecord(t *testing.T) {
	record := NormalizeRecord(Record{ID: "x", Name: "  A   Name ", Address: " Long Road ", Phone: " 1 ", Email: "A@EXAMPLE.COM", Tags: []string{"Wood", "wood"}, Status: StatusDraft, Version: 1})
	if record.Name != "A Name" || record.Email != "a@example.com" || len(record.Tags) != 1 {
		t.Fatalf("record=%+v", record)
	}
}

func TestTransitionRules(t *testing.T) {
	if err := ValidateTransition(StatusDraft, StatusApproved); err == nil {
		t.Fatal("expected invalid transition")
	}
	if err := ValidateTransition(StatusApproved, StatusArchived); err != nil {
		t.Fatal(err)
	}
}
