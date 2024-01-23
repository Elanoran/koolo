package action

import (
	"fmt"
	"strings"

	"github.com/Elanoran/d2go/pkg/data"
	"github.com/Elanoran/d2go/pkg/data/area"
	"github.com/Elanoran/d2go/pkg/data/item"
	"github.com/Elanoran/d2go/pkg/data/object"
	"github.com/Elanoran/d2go/pkg/data/stat"
	"github.com/Elanoran/d2go/pkg/itemfilter"
	"github.com/Elanoran/koolo/internal/action/step"
	"github.com/Elanoran/koolo/internal/config"
	stat2 "github.com/Elanoran/koolo/internal/event/stat"

	"github.com/Elanoran/koolo/internal/helper"
	"github.com/Elanoran/koolo/internal/hid"
	"github.com/Elanoran/koolo/internal/ui"
)

const (
	maxGoldPerStashTab = 2500000
	stashGoldBtnX      = 966
	stashGoldBtnY      = 526
)

func (b *Builder) Stash(forceStash bool) *Chain {
	return NewChain(func(d data.Data) (actions []Action) {
		b.logger.Debug("Checking for items to stash...")
		if !b.isStashingRequired(d, forceStash) {
			b.logger.Debug("No items to stash...")
			return
		}

		b.logger.Info("Stashing items...")

		switch d.PlayerUnit.Area {
		case area.KurastDocks:
			actions = append(actions, b.MoveToCoords(data.Position{X: 5146, Y: 5067}))
		case area.LutGholein:
			actions = append(actions, b.MoveToCoords(data.Position{X: 5130, Y: 5086}))
		}

		return append(actions,
			b.InteractObject(object.Bank,
				func(d data.Data) bool {
					return d.OpenMenus.Stash
				},
				step.SyncStep(func(d data.Data) error {
					b.stashGold(d)
					b.orderInventoryPotions(d)
					b.stashInventory(d, forceStash)
					hid.PressKey("esc")
					return nil
				}),
			),
		)
	})
}

func (b *Builder) orderInventoryPotions(d data.Data) {
	for _, i := range d.Items.ByLocation(item.LocationInventory) {
		if i.IsPotion() {
			if config.Config.Inventory.InventoryLock[i.Position.Y][i.Position.X] == 0 {
				continue
			}
			screenPos := ui.GetScreenCoordsForItem(i)
			hid.MovePointer(screenPos.X, screenPos.Y)
			helper.Sleep(100)
			hid.Click(hid.RightButton)
			helper.Sleep(200)
		}
	}
}

func (b *Builder) isStashingRequired(d data.Data, forceStash bool) bool {
	for _, i := range d.Items.ByLocation(item.LocationInventory) {
		if b.shouldStashIt(i, forceStash) {
			return true
		}
	}

	if d.PlayerUnit.Stats[stat.Gold] > d.PlayerUnit.MaxGold()/3 {
		return true
	}

	return false
}

func (b *Builder) stashGold(d data.Data) {
	gold, found := d.PlayerUnit.Stats[stat.Gold]
	if !found || gold == 0 {
		return
	}

	if d.PlayerUnit.Stats[stat.StashGold] < maxGoldPerStashTab {
		switchTab(1)
		clickStashGoldBtn()
		helper.Sleep(500)
	}

	for i := 2; i < 5; i++ {
		d = b.gr.GetData(false)
		gold, found = d.PlayerUnit.Stats[stat.Gold]
		if !found || gold == 0 {
			return
		}

		switchTab(i)
		clickStashGoldBtn()
	}
	b.logger.Info("All stash tabs are full of gold :D")
}

func (b *Builder) logStashItem(i data.Item) {
	// Check if the item has any stats
	if len(i.Stats) == 0 {
		b.logger.Debug("Item has no stats")
		return
	}

	// Initialize an empty slice to store the stats information
	var statsInfo []string

	// Iterate through each stat in the item's stats
	for statID, statData := range i.Stats {
		// Construct a string for each stat
		statInfo := fmt.Sprintf("%s %d", statID, statData.Value)

		// Add the statInfo to the statsInfo slice
		statsInfo = append(statsInfo, statInfo)
	}

	// Join the statsInfo slice into a single string
	statsString := strings.Join(statsInfo, ", ")

	// Construct the final textData
	textData := fmt.Sprintf("Item %s [%d] stashed - Stats: %s", i.Name, i.Quality, statsString)

	// Insert item into DB Log
	err := helper.DbLogText(textData)
	if err != nil {
		b.logger.Debug(fmt.Sprintf("%s", err))
	}
}

func (b *Builder) storeStashItem(i data.Item) {
	// Insert item into DB
	err := helper.DbLogItem(i)
	if err != nil {
		b.logger.Debug(fmt.Sprintf("%s", err))
	}
}

func (b *Builder) stashInventory(d data.Data, forceStash bool) {
	currentTab := 1
	switchTab(currentTab)

	for _, i := range d.Items.ByLocation(item.LocationInventory) {
		if !b.shouldStashIt(i, forceStash) {
			continue
		}
		for currentTab < 5 {
			if b.stashItemAction(i, forceStash) {
				b.logger.Debug(fmt.Sprintf("Item %s [%d] stashed", i.Name, i.Quality))
				// Log the item's stats along with its name and quality
				b.logStashItem(i)
				b.storeStashItem(i)
				break
			}
			if currentTab == 5 {
				// TODO: Stop the bot, stash is full
			}
			b.logger.Debug(fmt.Sprintf("Tab %d is full, switching to next one", currentTab))
			currentTab++
			switchTab(currentTab)
		}
	}
}

func (b *Builder) shouldStashIt(i data.Item, forceStash bool) bool {
	// Don't stash items from quests during leveling process, it makes things easier to track
	if _, isLevelingChar := b.ch.(LevelingCharacter); isLevelingChar && i.IsFromQuest() {
		return false
	}

	// Don't stash the Tomes, keys and WirtsLeg
	if i.Name == item.TomeOfTownPortal || i.Name == item.TomeOfIdentify || i.Name == item.Key || i.Name == "WirtsLeg" {
		return false
	}

	if config.Config.Inventory.InventoryLock[i.Position.Y][i.Position.X] == 0 || i.IsPotion() {
		return false
	}

	return forceStash || itemfilter.Evaluate(i, config.Config.Runtime.Rules)
}

func (b *Builder) stashItemAction(i data.Item, forceStash bool) bool {
	screenPos := ui.GetScreenCoordsForItem(i)
	hid.MovePointer(screenPos.X, screenPos.Y)
	helper.Sleep(170)
	screenshot := helper.Screenshot()
	helper.Sleep(150)
	hid.KeyDown("control")
	helper.Sleep(700)
	hid.Click(hid.LeftButton)
	helper.Sleep(700)
	hid.KeyUp("control")
	helper.Sleep(500)

	d := b.gr.GetData(false)
	for _, it := range d.Items.ByLocation(item.LocationInventory) {
		if it.UnitID == i.UnitID {
			return false
		}
	}

	// Don't log items that we already have in inventory during first run
	if !forceStash {
		stat2.ItemStashed(i, screenshot)
	}
	return true
}

func clickStashGoldBtn() {
	hid.MovePointer(stashGoldBtnX, stashGoldBtnY)
	helper.Sleep(170)
	hid.Click(hid.LeftButton)
	helper.Sleep(500)
	hid.PressKey("enter")
	helper.Sleep(700)
}

func switchTab(tab int) {
	x := 107
	y := 128
	tabSize := 82
	x = x + tabSize*tab - tabSize/2

	hid.MovePointer(x, y)
	helper.Sleep(400)
	hid.Click(hid.LeftButton)
	helper.Sleep(500)
}
