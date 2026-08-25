package model

import "fmt"

func NewAudit(recordID, action, actor, detail string, version int, sequence int64) AuditEvent {
	return AuditEvent{ID: fmt.Sprintf("audit-%06d", sequence), RecordID: recordID, Action: action, Actor: actor, Detail: detail, Version: version, CreatedAt: sequence}
}

func EventDetail(event AuditEvent) string {
	if event.Detail != "" {
		return event.Action + ": " + event.Detail
	}
	return event.Action
}

func NewWorkflow(recordID, name, owner string, steps []string) Workflow {
	copySteps := append([]string(nil), steps...)
	return Workflow{ID: recordID + "-workflow", RecordID: recordID, Name: name, State: "active", Steps: copySteps, Current: 0, Owner: owner, Revision: 1}
}

func (w Workflow) Complete() bool { return w.Current >= len(w.Steps) }

func (w Workflow) Advance(step string) (Workflow, error) {
	if w.Complete() {
		return w, fmt.Errorf("workflow is complete")
	}
	if w.Steps[w.Current] != step {
		return w, fmt.Errorf("expected step %q", w.Steps[w.Current])
	}
	w.Current++
	w.Revision++
	if w.Complete() {
		w.State = "complete"
	}
	return w, nil
}
