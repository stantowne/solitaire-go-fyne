package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// updateUI rebuilds the entire board
func (g *Game) updateUI() {
	// Update status
	stockCnt := len(g.stock)
	wasteCnt := len(g.waste)
	fCnt, tCnt := 0, 0
	for _, f := range g.foundations {
		fCnt += len(f)
	}
	for _, t := range g.tableau {
		tCnt += len(t)
	}
	g.statusLabel.SetText(fmt.Sprintf(
		"Moves: %d   Stock: %d   Waste: %d   Foundation: %d   Tableau: %d",
		g.moveCount, stockCnt, wasteCnt, fCnt, tCnt,
	))

	// Controls row with outlined buttons
	newGame := widget.NewButton("New Game", func() {
		g.resetPiles()
		g.shuffleAndDeal()
		g.updateUI()
	})
	togText := "Unmask Face Down Cards"
	if g.showDown {
		togText = "Mask Face Down Cards"
	}
	tog := widget.NewButton(togText, func() {
		g.toggleShow()
		g.updateUI()
	})
	newGameBox := container.NewStack(
		canvas.NewRectangle(color.NRGBA{R: 220, G: 220, B: 220, A: 255}),
		newGame,
	)
	newGameBox.Resize(fyne.NewSize(140, 40))
	togBox := container.NewStack(
		canvas.NewRectangle(color.NRGBA{R: 220, G: 220, B: 220, A: 255}),
		tog,
	)
	togBox.Resize(fyne.NewSize(200, 40))
	controls := container.NewHBox(newGameBox, togBox, layout.NewSpacer(), g.statusLabel)

	// --- Build grid cells ---

	// Left margin (adjust width as needed, e.g., 20px)
	leftMargin := canvas.NewRectangle(color.Transparent)
	leftMargin.SetMinSize(fyne.NewSize(20, cardSize.Height))

	// Stock
	var stockWidget fyne.CanvasObject
	if stockCnt == 0 {
		stockBtn := widget.NewButton("", func() {
			g.drawCards()
			g.updateUI()
		})
		stockBtn.Importance = widget.LowImportance
		stockBtn.Resize(cardSize)
		wrapper := container.NewWithoutLayout()
		wrapper.Add(stockBtn)
		wrapper.Resize(cardSize)
		stockWidget = wrapper
	} else {
		stockBtn := widget.NewButton("", func() {
			g.drawCards()
			g.updateUI()
		})
		stockBtn.Importance = widget.LowImportance
		stockBtn.Resize(cardSize)
		stockRect := canvas.NewRectangle(cardBackColor)
		stockRect.SetMinSize(cardSize)
		stockRect.StrokeColor = color.Black
		stockRect.StrokeWidth = 1
		stock := container.NewStack(stockRect, stockBtn)
		stock.Resize(cardSize)
		wrapper := container.NewWithoutLayout()
		wrapper.Add(stock)
		wrapper.Resize(cardSize)
		stockWidget = wrapper
	}

	// Waste (double width visually, but occupies one cell)
	var wasteWidget fyne.CanvasObject
	wasteWidth := cardSize.Width*2 + 10 // visually double width
	if len(g.waste) == 0 {
		r := canvas.NewRectangle(color.Transparent)
		r.SetMinSize(fyne.NewSize(wasteWidth, cardSize.Height))
		wasteStack := container.NewStack(r)
		wasteStack.Resize(fyne.NewSize(wasteWidth, cardSize.Height))
		wasteWidget = wasteStack
	} else {
		renderedWaste := g.renderWaste(g.waste)
		wasteWidget = container.NewStack(renderedWaste)
		wasteWidget.Resize(fyne.NewSize(wasteWidth, cardSize.Height))
	}

	// Empty cell (col 3)
	emptyCell := canvas.NewRectangle(color.Transparent)
	emptyCell.SetMinSize(cardSize)

	// Foundations
	fWidgets := make([]fyne.CanvasObject, 4)
	for i := 0; i < 4; i++ {
		fWidgets[i] = g.renderPileOrEmpty(g.foundations[i], "foundation", i)
	}

	// Tableau
	tCols := make([]fyne.CanvasObject, 7)
	for i := 0; i < 7; i++ {
		tCols[i] = g.renderPileOrEmpty(g.tableau[i], "tableau", i)
	}

	// Build the grid: 2 rows × 8 columns (flat slice)
	// Row 0: left margin, stock, waste, empty, foundations (4)
	// Row 1: left margin, tableau (7)
	gridCells := []fyne.CanvasObject{
		stockWidget, wasteWidget, emptyCell, fWidgets[0], fWidgets[1], fWidgets[2], fWidgets[3],
		tCols[0], tCols[1], tCols[2], tCols[3], tCols[4], tCols[5], tCols[6],
	}

	// Set up the grid
	grid := container.New(layout.NewGridWrapLayout(cardSize), gridCells...)

	// Add vertical spacer between controls and grid
	verticalSpacer := canvas.NewRectangle(color.Transparent)
	verticalSpacer.SetMinSize(fyne.NewSize(0, 30))

	g.board.Objects = []fyne.CanvasObject{
		controls,
		verticalSpacer,
		grid,
	}
	g.board.Refresh()
}

// renderPileOrEmpty shows a gray placeholder or stacks cards
func (g *Game) renderPileOrEmpty(pile []*Card, pileType string, pileIndex int) fyne.CanvasObject {
	if len(pile) == 0 {
		return NewEmptyPileWidget(g, pileType, pileIndex)
	}

	if pileType == "foundation" {
		return g.renderFoundation(pile)
	}
	return g.renderPile(pile)
}

// renderPile stacks every card vertically with spacing
func (g *Game) renderPile(pile []*Card) *fyne.Container {
	stack := container.NewWithoutLayout()
	for idx, c := range pile {
		w := NewCardWidget(c, g)
		w.Resize(cardSize)
		y := float32(idx) * cardOverlap
		w.Move(fyne.NewPos(0, y))
		stack.Add(w)
	}
	totalH := float32(len(pile)-1)*cardOverlap + cardSize.Height
	stack.Resize(fyne.NewSize(cardSize.Width, totalH))
	return stack
}

// renderWaste shows up to 3 cards fanned horizontally
func (g *Game) renderWaste(pile []*Card) *fyne.Container {
	cont := container.NewWithoutLayout()
	start := 0
	if len(pile) > 3 {
		start = len(pile) - 3
	}
	visible := pile[start:]
	for i, c := range visible {
		w := NewCardWidget(c, g)
		w.Resize(cardSize)
		w.Move(fyne.NewPos(float32(i)*wasteOffsetX, 0))
		cont.Add(w)
	}
	width := float32(len(visible)-1)*wasteOffsetX + cardSize.Width
	cont.Resize(fyne.NewSize(width, cardSize.Height))
	return cont
}

func (g *Game) renderFoundation(pile []*Card) *fyne.Container {
	stack := container.NewWithoutLayout()
	if len(pile) > 0 {
		top := pile[len(pile)-1]
		w := NewCardWidget(top, g)
		w.Resize(cardSize)
		w.Move(fyne.NewPos(0, 0))
		stack.Add(w)
	}
	stack.Resize(cardSize)
	return stack
}
