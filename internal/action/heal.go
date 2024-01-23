package action

import (
	"fmt"

	"github.com/Elanoran/d2go/pkg/data"
	"github.com/Elanoran/koolo/internal/action/step"
	"github.com/Elanoran/koolo/internal/helper"
	"github.com/Elanoran/koolo/internal/hid"
	"github.com/Elanoran/koolo/internal/town"
)

func (b *Builder) Heal() *Chain {
	return NewChain(func(d data.Data) []Action {
		shouldHeal := false
		if d.PlayerUnit.HPPercent() < 80 {
			b.logger.Info(fmt.Sprintf("Current life is %d%%, healing on NPC", d.PlayerUnit.HPPercent()))
			shouldHeal = true
		}

		if d.PlayerUnit.HasDebuff() {
			b.logger.Info(fmt.Sprintf("Debuff detected, healing on NPC"))
			shouldHeal = true
		}

		if shouldHeal {
			return []Action{b.InteractNPC(
				town.GetTownByArea(d.PlayerUnit.Area).HealNPC(),
				step.SyncStep(func(d data.Data) error {
					helper.Sleep(300)
					hid.PressKey("esc")
					helper.Sleep(100)
					return nil
				}),
			)}
		}

		return nil
	})
}
