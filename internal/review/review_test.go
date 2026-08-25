package review

import (
	"aromaatelier/internal/model"
	"aromaatelier/internal/store"
	"testing"
)

func TestReviewLifecycle(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/review.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.CreateRecord(model.Record{ID: "r2", Name: "Rose", Address: "No. 2 Lane", Phone: "101"}); err != nil {
		t.Fatal(err)
	}
	s := New(db)
	if _, err := s.Submit("r2", "a"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Approve("r2", "b", "looks good"); err != nil {
		t.Fatal(err)
	}
	got, _ := db.GetRecord("r2")
	if got.Status != model.StatusApproved {
		t.Fatalf("status=%s", got.Status)
	}
}
