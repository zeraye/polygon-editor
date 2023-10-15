package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	polygon_editor "github.com/zeraye/polygon-editor"
	"github.com/zeraye/polygon-editor/pkg/config"
)

func main() {
	config, err := config.LoadStandard("config", "config.toml")
	if err != nil {
		log.Fatal(err)
	}

	app := app.New()
	window := app.NewWindow(config.Window.Name)
	game := polygon_editor.NewGame(config)

	window.SetContent(game.BuildUI())
	window.Resize(fyne.NewSize(float32(config.Window.Width), float32(config.Window.Height)))
	window.SetFixedSize(config.Window.FixedSize)
	window.ShowAndRun()
}
