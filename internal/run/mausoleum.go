package run

import (
	"github.com/hectorgimenez/d2go/pkg/data"
	"github.com/hectorgimenez/d2go/pkg/data/area"
	"github.com/hectorgimenez/koolo/internal/action"
	"github.com/hectorgimenez/koolo/internal/config"
	"github.com/hectorgimenez/koolo/internal/health"
)

type Mausoleum struct {
	baseRun
	bm health.BeltManager
}

func (a Mausoleum) Name() string {
	return "Mausoleum"
}

func (a Mausoleum) BuildActions() (actions []action.Action) {
	actions = append(actions,
		a.builder.WayPoint(area.ColdPlains),
		a.builder.MoveToArea(area.BurialGrounds),
		a.builder.MoveToArea(area.Mausoleum),
		a.builder.ClearArea(true, data.MonsterAnyFilter()),
	)

	// Go back to town to buy potions if needed
	actions = append(actions, action.NewChain(func(d data.Data) []action.Action {
		if config.Config.Character.BuyPots && (a.bm.ShouldBuyPotions(d)) {
			return a.builder.InRunReturnTownRoutine()
		}

		return nil
	}))

	return
}
