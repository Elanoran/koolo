package helper

import (
	"github.com/Elanoran/d2go/pkg/data"
	"github.com/Elanoran/d2go/pkg/data/area"
	"github.com/Elanoran/d2go/pkg/data/skill"
	"github.com/Elanoran/koolo/internal/config"
)

func CanTeleport(d data.Data) bool {
	_, found := d.PlayerUnit.Skills[skill.Teleport]

	// Duriel's Lair is bugged and teleport doesn't work here
	if d.PlayerUnit.Area == area.DurielsLair {
		return false
	}

	return found && config.Config.Bindings.Teleport != "" && !d.PlayerUnit.Area.IsTown()
}
