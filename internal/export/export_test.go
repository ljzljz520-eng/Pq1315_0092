package export

import (
	"aromaatelier/internal/model"
	"aromaatelier/internal/store"
	"strings"
	"testing"
)

func TestJSONUsesCurrentAddress(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/export.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.CreateRecord(model.Record{ID: "r3", Name: "Mint", Address: "Old Street", Phone: "102"}); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveRecord(model.Record{ID: "r3", Name: "Mint", Address: "New Street", Phone: "102", Status: model.StatusDraft, Version: 2}); err != nil {
		t.Fatal(err)
	}
	payload, err := New(db).JSON("r3")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(payload), "New Street") {
		t.Fatal("new address absent")
	}
}
