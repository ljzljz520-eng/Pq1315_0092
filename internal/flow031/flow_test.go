package flow031

import (
	"aromaatelier/internal/model"
	"aromaatelier/internal/store"
	"strings"
	"testing"
)

func testDB(t *testing.T) *store.DB {
	t.Helper()
	db, err := store.Open(t.TempDir() + "/flow.db")
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestWorkflowCreateReviewArchive(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	h := New(db)
	if _, err := h.Register(model.Record{ID: "wf1", Name: "Candle", Address: "Studio Road", Phone: "1"}, "maker"); err != nil {
		t.Fatal(err)
	}
	if _, err := h.Submit("wf1", "maker"); err != nil {
		t.Fatal(err)
	}
	if _, err := h.Approve("wf1", "reviewer", "approved"); err != nil {
		t.Fatal(err)
	}
	if _, err := h.Archive("wf1", "maker"); err != nil {
		t.Fatal(err)
	}
	got, _ := h.Search(model.SearchQuery{Text: "Candle"})
	if len(got) != 1 || got[0].Status != model.StatusArchived {
		t.Fatalf("got=%v", got)
	}
}

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	h := New(db)
	if _, err := h.Register(model.Record{ID: "wf2", Name: "Moss", Address: "North Studio", Phone: "2"}, "maker"); err != nil {
		t.Fatal(err)
	}
	if _, err := h.UpdateAddress("wf2", "maker", "South Studio", "relocated"); err != nil {
		t.Fatal(err)
	}
	if _, err := h.Submit("wf2", "maker"); err != nil {
		t.Fatal(err)
	}
	if _, err := h.Approve("wf2", "reviewer", "ready"); err != nil {
		t.Fatal(err)
	}
	payload, err := h.ExportOne("wf2")
	if err != nil || !strings.Contains(string(payload), "South Studio") {
		t.Fatalf("payload=%s err=%v", payload, err)
	}
}

func TestWorkflowImportReport(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	h := New(db)
	input := "id\tname\taddress\tphone\temail\ttags\ni1\tIncense\tWest Studio\t3\ti@example.com\twood"
	report := h.ImportAndReport(input, "importer")
	if report.Created != 1 || len(report.Published) != 1 {
		t.Fatalf("report=%+v", report)
	}
}

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := t.TempDir() + "/persist.db"
	db, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	h := New(db)
	if _, err := h.Register(model.Record{ID: "persist", Name: "Sandal", Address: "Old Plaza", Phone: "4"}, "maker"); err != nil {
		t.Fatal(err)
	}
	if _, err := h.UpdateAddress("persist", "maker", "New Plaza", "move"); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	record, err := reopened.GetRecord("persist")
	if err != nil {
		t.Fatal(err)
	}
	if record.Address != "New Plaza" {
		t.Fatalf("address=%s", record.Address)
	}
}

func TestSearchReturnsStableOrder(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	h := New(db)
	for _, id := range []string{"b", "a"} {
		if _, err := h.Register(model.Record{ID: id, Name: id, Address: "Shared Lane", Phone: id}, "x"); err != nil {
			t.Fatal(err)
		}
	}
	rows, err := h.Search(model.SearchQuery{Text: "shared"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[0].ID != "a" {
		t.Fatalf("rows=%v", rows)
	}
}

func Test1315BusinessRegression(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	h := New(db)
	if _, err := h.Register(model.Record{ID: "brand-1", Name: "First", Address: "First Old", Phone: "10"}, "maker"); err != nil {
		t.Fatal(err)
	}
	if _, err := h.Register(model.Record{ID: "brand-2", Name: "Second", Address: "Second Old", Phone: "11"}, "maker"); err != nil {
		t.Fatal(err)
	}
	if _, err := h.UpdateAddress("brand-2", "maker", "Second Current", "new workshop"); err != nil {
		t.Fatal(err)
	}
	exports, err := h.ExportSearch(model.SearchQuery{Text: ""})
	if err != nil {
		t.Fatal(err)
	}
	if len(exports) != 2 || !strings.Contains(string(exports[1]), "Second Current") {
		t.Fatalf("exports=%s", exports[1])
	}
}
