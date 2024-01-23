package action

import (
	"github.com/Elanoran/koolo/internal/health"
	"github.com/Elanoran/koolo/internal/reader"
	"github.com/Elanoran/koolo/internal/town"
	"github.com/Elanoran/koolo/internal/ui"
	"go.uber.org/zap"
)

type Builder struct {
	logger *zap.Logger
	sm     town.ShopManager
	bm     health.BeltManager
	gr     *reader.GameReader
	ch     Character
	tf     *ui.TemplateFinder
}

func NewBuilder(logger *zap.Logger, sm town.ShopManager, bm health.BeltManager, gr *reader.GameReader, ch Character, tf *ui.TemplateFinder) *Builder {
	return &Builder{
		logger: logger,
		sm:     sm,
		bm:     bm,
		gr:     gr,
		ch:     ch,
		tf:     tf,
	}
}
