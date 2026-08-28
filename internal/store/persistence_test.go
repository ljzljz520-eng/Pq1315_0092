package store

import (
	"aromaatelier/internal/model"
	"testing"
)

func TestRecordVersionConflict(t *testing.T) {
	db, err := Open(t.TempDir() + "/db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.CreateRecord(model.Record{ID: "x", Name: "X", Address: "Some Lane", Phone: "1"}); err != nil {
		t.Fatal(err)
	}
	err = db.SaveRecord(model.Record{ID: "x", Name: "X", Address: "Other Lane", Phone: "1", Status: model.StatusDraft, Version: 3})
	if err != ErrConflict {
		t.Fatalf("err=%v", err)
	}
}

func TestAttachmentChecksum(t *testing.T) {
	db, err := Open(t.TempDir() + "/db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	err = db.SaveAttachment(model.Attachment{ID: "a", RecordID: "r", Name: "label", Content: "hello"})
	if err != nil {
		t.Fatal(err)
	}
	attachment, err := db.GetAttachment("a")
	if err != nil || attachment.Checksum == "" {
		t.Fatalf("attachment=%+v err=%v", attachment, err)
	}
}
