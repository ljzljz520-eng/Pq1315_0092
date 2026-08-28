package export

import (
	"encoding/json"
	"fmt"
	"sort"

	"aromaatelier/internal/model"
	"aromaatelier/internal/store"
)

type Service struct{ store *store.DB }

func New(db *store.DB) *Service { return &Service{store: db} }

func (s *Service) Document(id string) (model.ExportDocument, error) {
	record, err := s.store.GetRecord(id)
	if err != nil {
		return model.ExportDocument{}, err
	}
	audit, err := s.store.AuditFor(id)
	if err != nil {
		return model.ExportDocument{}, err
	}
	sort.Slice(audit, func(i, j int) bool { return audit[i].CreatedAt < audit[j].CreatedAt })
	return model.ExportDocument{RecordID: record.ID, Name: record.Name, Address: record.Address, Phone: record.Phone, Status: string(record.Status), Version: record.Version, Audit: audit}, nil
}

func (s *Service) JSON(id string) ([]byte, error) {
	doc, err := s.Document(id)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(doc, "", "  ")
}

func (s *Service) CSV(query model.SearchQuery) (string, error) {
	rows, err := s.store.ListRecords(query)
	if err != nil {
		return "", err
	}
	result := "id,name,address,phone,status,version\n"
	for _, row := range rows {
		result += fmt.Sprintf("%s,%s,%s,%s,%s,%d\n", row.ID, row.Name, row.Address, row.Phone, row.Status, row.Version)
	}
	return result, nil
}
