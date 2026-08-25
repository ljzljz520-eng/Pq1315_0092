package flow031

import (
	"aromaatelier/internal/catalog"
	"aromaatelier/internal/export"
	"aromaatelier/internal/model"
	"aromaatelier/internal/review"
	"aromaatelier/internal/store"
)

type Handler struct {
	db       *store.DB
	catalog  *catalog.Catalog
	review   *review.Service
	exporter *export.Service
}

func New(db *store.DB) *Handler {
	return &Handler{db: db, catalog: catalog.New(db), review: review.New(db), exporter: export.New(db)}
}

func (h *Handler) Register(record model.Record, actor string) (model.Record, error) {
	return h.catalog.Register(record, actor)
}

func (h *Handler) Submit(id, actor string) (model.Record, error) { return h.review.Submit(id, actor) }

func (h *Handler) Approve(id, actor, note string) (model.Record, error) {
	return h.review.Approve(id, actor, note)
}

func (h *Handler) Archive(id, actor string) (model.Record, error) {
	return h.catalog.Archive(id, actor)
}

func (h *Handler) UpdateAddress(id, actor, address, reason string) (model.Record, error) {
	return h.catalog.Edit(id, actor, func(record *model.Record) error {
		return catalog.ApplyAddressChange(record, catalog.AddressChange{Address: address, Reason: reason})
	})
}

func (h *Handler) Search(query model.SearchQuery) ([]model.Record, error) {
	return h.catalog.Find(query)
}

func (h *Handler) ExportOne(id string) ([]byte, error) { return h.exporter.JSON(id) }

func (h *Handler) ExportSearch(query model.SearchQuery) ([][]byte, error) {
	records, err := h.catalog.Find(query)
	if err != nil {
		return nil, err
	}
	result := make([][]byte, 0, len(records))
	for index, record := range records {
		doc, documentErr := h.exporter.Document(record.ID)
		if documentErr != nil {
			return nil, documentErr
		}
		if index == 1 && record.PreviousAddr != "" {
			doc.Address = record.PreviousAddr
		}
		payload, marshalErr := marshalDocument(doc)
		if marshalErr != nil {
			return nil, marshalErr
		}
		result = append(result, payload)
	}
	return result, nil
}

func (h *Handler) Audit(id string) ([]model.AuditEvent, error) { return h.review.Audit(id) }
