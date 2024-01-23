package action

import (
	"fmt"

	"github.com/Elanoran/d2go/pkg/data"
	"github.com/Elanoran/d2go/pkg/data/item"
	"github.com/Elanoran/d2go/pkg/data/stat"
	"github.com/Elanoran/koolo/internal/action/step"
	"github.com/Elanoran/koolo/internal/config"
	"github.com/Elanoran/koolo/internal/helper"
	"github.com/Elanoran/koolo/internal/hid"
	"github.com/Elanoran/koolo/internal/ui"
)

func (b *Builder) IdentifyAll(skipIdentify bool) *Chain {
	return NewChain(func(d data.Data) (actions []Action) {
		items := b.itemsToIdentify(d)

		b.logger.Debug("Checking for items to identify...")
		if len(items) == 0 || skipIdentify {
			b.logger.Debug("No items to identify...")
			return
		}

		helper.Sleep(2000)

		idTome, found := d.Items.Find(item.TomeOfIdentify, item.LocationInventory)
		if !found {
			b.logger.Warn("ID Tome not found, not identifying items")
			return
		}

		if st, statFound := idTome.Stats[stat.Quantity]; !statFound || st.Value < len(items) {
			b.logger.Info("Not enough ID scrolls, refilling...")
			actions = append(actions, b.VendorRefill(true, false))
		}

		b.logger.Info(fmt.Sprintf("Identifying %d items...", len(items)))
		actions = append(actions, NewStepChain(func(d data.Data) []step.Step {
			return []step.Step{
				step.SyncStepWithCheck(func(d data.Data) error {
					hid.PressKey(config.Config.Bindings.OpenInventory)
					return nil
				}, func(d data.Data) step.Status {
					if d.OpenMenus.Inventory {
						return step.StatusCompleted
					}
					return step.StatusInProgress
				}),
				step.SyncStep(func(d data.Data) error {

					for _, i := range items {
						identifyItem(idTome, i)
					}

					hid.PressKey("esc")

					return nil
				}),
			}
		}))

		return
	}, Resettable(), CanBeSkipped())
}

func (b *Builder) itemsToIdentify(d data.Data) (items []data.Item) {
	for _, i := range d.Items.ByLocation(item.LocationInventory) {
		if i.Identified || i.Quality == item.QualityNormal || i.Quality == item.QualitySuperior {
			continue
		}

		items = append(items, i)
	}

	return
}

func identifyItem(idTome data.Item, i data.Item) {
	screenPos := ui.GetScreenCoordsForItem(idTome)
	hid.MovePointer(screenPos.X, screenPos.Y)
	helper.Sleep(500)
	hid.Click(hid.RightButton)
	helper.Sleep(1000)

	screenPos = ui.GetScreenCoordsForItem(i)
	hid.MovePointer(screenPos.X, screenPos.Y)
	helper.Sleep(1000)
	hid.Click(hid.LeftButton)
	helper.Sleep(350)
}
