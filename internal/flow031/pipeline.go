package flow031

import (
	"aromaatelier/internal/importer"
	"aromaatelier/internal/model"
	"fmt"
)

type ImportReport struct {
	importer.Result
	Published []string
}

func (h *Handler) ImportAndReport(input string, actor string) ImportReport {
	rows := importer.ParseTSV(input)
	result := importer.New(h.catalog).Import(rows, actor)
	report := ImportReport{Result: result, Published: []string{}}
	for _, row := range rows {
		if _, err := h.Submit(row.ExternalID, actor); err != nil {
			continue
		}
		if _, err := h.Approve(row.ExternalID, actor, "import review"); err != nil {
			continue
		}
		report.Published = append(report.Published, row.ExternalID)
	}
	return report
}

func (h *Handler) CreateWorkflow(id, owner string) (model.Workflow, error) {
	if _, err := h.db.GetRecord(id); err != nil {
		return model.Workflow{}, err
	}
	w := model.NewWorkflow(id, "brand publication", owner, []string{"register", "review", "confirm", "archive"})
	if err := h.db.SaveWorkflow(w); err != nil {
		return model.Workflow{}, err
	}
	return w, nil
}

func (h *Handler) AdvanceWorkflow(id, step string) (model.Workflow, error) {
	w, err := h.db.GetWorkflow(id + "-workflow")
	if err != nil {
		return model.Workflow{}, err
	}
	next, err := w.Advance(step)
	if err != nil {
		return model.Workflow{}, fmt.Errorf("advance: %w", err)
	}
	if err := h.db.SaveWorkflow(next); err != nil {
		return model.Workflow{}, err
	}
	return next, nil
}
