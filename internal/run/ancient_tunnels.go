package run

import (
	"github.com/Elanoran/d2go/pkg/data"
	"github.com/Elanoran/d2go/pkg/data/area"
	"github.com/Elanoran/koolo/internal/action"
	"github.com/Elanoran/koolo/internal/action/step"
	"github.com/Elanoran/koolo/internal/config"
)

type AncientTunnels struct {
	baseRun
}

func (a AncientTunnels) Name() string {
	return "AncientTunnels"
}

func (a AncientTunnels) BuildActions() []action.Action {
	actions := []action.Action{
		a.builder.WayPoint(area.LostCity),         // Moving to starting point (Lost City)
		a.builder.MoveToArea(area.AncientTunnels), // Travel to ancient tunnels
	}

	if config.Config.Companion.Enabled && config.Config.Companion.Leader {
		actions = append(actions, action.NewStepChain(func(_ data.Data) []step.Step {
			return []step.Step{step.OpenPortal()}
		}))
	}

	// Clear Ancient Tunnels
	return append(actions, a.builder.ClearArea(true, data.MonsterAnyFilter()))
}
