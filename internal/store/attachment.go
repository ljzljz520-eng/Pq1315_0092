package store

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"aromaatelier/internal/model"
	"go.etcd.io/bbolt"
)

func (s *DB) SaveAttachment(attachment model.Attachment) error {
	if attachment.ID == "" || attachment.RecordID == "" || attachment.Name == "" {
		return fmt.Errorf("attachment identity is required")
	}
	if attachment.Checksum == "" {
		digest := sha256.Sum256([]byte(attachment.Content))
		attachment.Checksum = hex.EncodeToString(digest[:])
	}
	return s.transaction(func(tx *bbolt.Tx) error { return putJSON(tx.Bucket([]byte("attachments")), attachment.ID, attachment) })
}

func (s *DB) GetAttachment(id string) (model.Attachment, error) {
	var attachment model.Attachment
	err := s.view(func(tx *bbolt.Tx) error { return getJSON(tx.Bucket([]byte("attachments")), id, &attachment) })
	return attachment, err
}

func (s *DB) ListAttachments(recordID string) ([]model.Attachment, error) {
	var list []model.Attachment
	err := s.view(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("attachments")).ForEach(func(_, value []byte) error {
			if value == nil {
				return nil
			}
			var attachment model.Attachment
			if err := unmarshal(value, &attachment); err != nil {
				return err
			}
			if recordID == "" || attachment.RecordID == recordID {
				list = append(list, attachment)
			}
			return nil
		})
	})
	return list, err
}
