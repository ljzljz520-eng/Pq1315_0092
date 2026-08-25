package catalog

import (
	"fmt"
	"strings"

	"aromaatelier/internal/model"
	"aromaatelier/internal/store"
)

type Catalog struct {
	store *store.DB
}

func New(db *store.DB) *Catalog { return &Catalog{store: db} }

func (c *Catalog) Register(input model.Record, actor string) (model.Record, error) {
	input.ID = strings.TrimSpace(input.ID)
	input.CreatedBy = actor
	input.UpdatedBy = actor
	input.Status = model.StatusDraft
	input.Version = 1
	if err := c.store.CreateRecord(input); err != nil {
		return model.Record{}, err
	}
	if _, err := c.store.AppendAudit(input.ID, "registered", actor, "draft created", input.Version); err != nil {
		return model.Record{}, err
	}
	return c.store.GetRecord(input.ID)
}

func (c *Catalog) Find(query model.SearchQuery) ([]model.Record, error) {
	if query.Limit < 0 {
		return nil, fmt.Errorf("limit cannot be negative")
	}
	return c.store.ListRecords(query)
}

func (c *Catalog) Get(id string) (model.Record, error) { return c.store.GetRecord(id) }

func (c *Catalog) Edit(id, actor string, update func(*model.Record) error) (model.Record, error) {
	record, err := c.store.GetRecord(id)
	if err != nil {
		return model.Record{}, err
	}
	if record.Status == model.StatusArchived {
		return model.Record{}, fmt.Errorf("archived record cannot be edited")
	}
	if err := update(&record); err != nil {
		return model.Record{}, err
	}
	record.Version++
	record.UpdatedBy = actor
	if err := c.store.SaveRecord(record); err != nil {
		return model.Record{}, err
	}
	if _, err := c.store.AppendAudit(id, "updated", actor, "record details changed", record.Version); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (c *Catalog) Archive(id, actor string) (model.Record, error) {
	record, err := c.store.GetRecord(id)
	if err != nil {
		return model.Record{}, err
	}
	if err := model.ValidateTransition(record.Status, model.StatusArchived); err != nil {
		return model.Record{}, err
	}
	record.Status = model.StatusArchived
	record.Version++
	record.UpdatedBy = actor
	if err := c.store.SaveRecord(record); err != nil {
		return model.Record{}, err
	}
	if _, err := c.store.AppendAudit(id, "archived", actor, "record archived", record.Version); err != nil {
		return model.Record{}, err
	}
	return record, nil
}
