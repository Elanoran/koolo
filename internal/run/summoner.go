package run

import (
	"github.com/Elanoran/d2go/pkg/data"
	"github.com/Elanoran/d2go/pkg/data/area"
	"github.com/Elanoran/d2go/pkg/data/npc"
	"github.com/Elanoran/koolo/internal/action"
)

type Summoner struct {
	baseRun
}

func (s Summoner) Name() string {
	return "Summoner"
}

func (s Summoner) BuildActions() (actions []action.Action) {
	return []action.Action{
		s.builder.WayPoint(area.ArcaneSanctuary), // Moving to starting point (Arcane Sanctuary)
		s.builder.MoveTo(func(d data.Data) (data.Position, bool) {
			m, found := d.NPCs.FindOne(npc.Summoner)

			return m.Positions[0], found
		}), // Travel to boss position
		s.char.KillSummoner(), // Kill Summoner
	}
}
