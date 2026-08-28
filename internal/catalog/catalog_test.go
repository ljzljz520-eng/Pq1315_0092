package catalog

import (
	"aromaatelier/internal/model"
	"aromaatelier/internal/store"
	"testing"
)

func TestRegisterAndFind(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/catalog.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	c := New(db)
	_, err = c.Register(model.Record{ID: "r1", Name: "Cedar", Address: "No. 1 Lane", Phone: "100", Email: "a@example.com"}, "maker")
	if err != nil {
		t.Fatal(err)
	}
	rows, err := c.Find(model.SearchQuery{Text: "cedar"})
	if err != nil || len(rows) != 1 {
		t.Fatalf("rows=%v err=%v", rows, err)
	}
}
