package popover

import (
	"crypto/rand"
)

type TooltipProps struct {
	TriggerID   string
	TooltipID   string
	TargetClass string
}

func (t *TooltipProps) GetTriggerID() string {
	if t.TriggerID == "" {
		t.TriggerID = rand.Text()
	}
	return t.TriggerID
}

func (t *TooltipProps) GetTooltipID() string {
	if t.TooltipID == "" {
		t.TooltipID = rand.Text()
	}
	return t.TooltipID
}
