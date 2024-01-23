package main

import (
	"log"

	"github.com/Elanoran/d2go/pkg/memory"
	"github.com/Elanoran/koolo/internal/action"
	"github.com/Elanoran/koolo/internal/character"
	"github.com/Elanoran/koolo/internal/config"
	"github.com/Elanoran/koolo/internal/health"
	"github.com/Elanoran/koolo/internal/helper"
	"github.com/Elanoran/koolo/internal/hid"
	"github.com/Elanoran/koolo/internal/reader"
	"github.com/Elanoran/koolo/internal/town"
	"github.com/Elanoran/koolo/internal/ui"
	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
	"go.uber.org/zap"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err.Error())
	}

	logger, _ := zap.NewDevelopment()

	window := robotgo.FindWindow("Diablo II: Resurrected")
	if window == win.HWND_TOP {
		panic(err)
	}
	win.SetForegroundWindow(window)

	pos := win.WINDOWPLACEMENT{}
	point := win.POINT{}
	win.ClientToScreen(window, &point)
	win.GetWindowPlacement(window, &pos)

	hid.WindowLeftX = int(point.X + 1)
	hid.WindowTopY = int(point.Y)
	hid.GameAreaSizeX = int(pos.RcNormalPosition.Right) - hid.WindowLeftX - 9
	hid.GameAreaSizeY = int(pos.RcNormalPosition.Bottom) - hid.WindowTopY - 9
	helper.Sleep(1000)

	process, err := memory.NewProcess()
	if err != nil {
		panic(err)
	}

	gd := memory.NewGameReader(process)
	gr := &reader.GameReader{GameReader: gd}
	bm := health.NewBeltManager(logger)
	sm := town.NewShopManager(logger, bm)
	char, err := character.BuildCharacter(logger)
	if err != nil {
		panic(err)
	}
	tf, err := ui.NewTemplateFinder(logger, "assets")
	if err != nil {
		panic(err)
	}

	b := action.NewBuilder(logger, sm, bm, gr, char, tf)
	a := b.ItemPickup(true, -1)

	gr.GetData(true)

	for err == nil {
		data := gr.GetData(false)
		err = a.NextStep(logger, data)
	}
}
