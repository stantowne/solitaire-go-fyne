package main

/*
Klondike Solitaire - Go/Fyne Implementation

*/

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// main sets up and runs the Fyne window
func main() {
	fmt.Println("Starting Klondike Solitaire...")
	a := app.New()
	w := a.NewWindow("Klondike Solitaire")

	game := NewGame()
	game.statusLabel = widget.NewLabel("")
	game.board = container.NewVBox()
	w.SetContent(game.board)

	game.updateUI()
	w.Resize(fyne.NewSize(900, 600))
	w.ShowAndRun()
}
