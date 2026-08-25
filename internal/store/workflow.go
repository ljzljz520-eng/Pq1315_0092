package store

import (
	"aromaatelier/internal/model"
	"go.etcd.io/bbolt"
)

func (s *DB) SaveWorkflow(workflow model.Workflow) error {
	if workflow.ID == "" || workflow.RecordID == "" {
		return ErrNotFound
	}
	return s.transaction(func(tx *bbolt.Tx) error { return putJSON(tx.Bucket([]byte("workflows")), workflow.ID, workflow) })
}

func (s *DB) GetWorkflow(id string) (model.Workflow, error) {
	var workflow model.Workflow
	err := s.view(func(tx *bbolt.Tx) error { return getJSON(tx.Bucket([]byte("workflows")), id, &workflow) })
	return workflow, err
}

func (s *DB) ListWorkflows(recordID string) ([]model.Workflow, error) {
	var list []model.Workflow
	err := s.view(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("workflows")).ForEach(func(_, value []byte) error {
			if value == nil {
				return nil
			}
			var workflow model.Workflow
			if err := unmarshal(value, &workflow); err != nil {
				return err
			}
			if recordID == "" || workflow.RecordID == recordID {
				list = append(list, workflow)
			}
			return nil
		})
	})
	return list, err
}
