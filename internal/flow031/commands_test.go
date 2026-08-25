package flow031

import (
	"aromaatelier/internal/model"
	"testing"
)

func TestPublishChecklist(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	h := New(db)
	if _, err := h.Register(model.Record{ID: "pub", Name: "Publishable", Address: "A Long Studio", Phone: "5"}, "maker"); err != nil {
		t.Fatal(err)
	}
	if _, err := h.Submit("pub", "maker"); err != nil {
		t.Fatal(err)
	}
	result, err := h.Publish("pub", "reviewer", model.ReviewChecklist{IdentityVerified: true, AddressVerified: true, SafetyNote: true, PhotoAttached: true, DescriptionReady: true})
	if err != nil || result.Status != "approved" {
		t.Fatalf("result=%+v err=%v", result, err)
	}
}
