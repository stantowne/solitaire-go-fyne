package main

import (
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// Game holds all piles and UI state
type Game struct {
	stock       []*Card
	waste       []*Card
	foundations [4][]*Card
	tableau     [7][]*Card

	moveCount int
	showDown  bool
	drawCount int

	board       *fyne.Container
	statusLabel *widget.Label

	// card selected on first click
	firstSelectedCard   *Card
	firstSelectedSource string // "waste", "foundation", "tableau"
	firstSelectedIndex  int    // which pile (e.g., tableau column or foundation index)
	firstSelectedDepth  int    // for tableau: index in the slice (for moving groups)

	// card selected on second click
	secondSelectedCard   *Card
	secondSelectedSource string // "waste", "foundation", "tableau"
	secondSelectedIndex  int    // which pile (e.g., tableau column or foundation index)

}

// NewGame creates, shuffles, and deals a fresh layout
func NewGame() *Game {
	g := &Game{drawCount: 3}
	g.resetPiles()
	g.shuffleAndDeal()
	return g
}

func (g *Game) cancelFirstSelection() {
	g.firstSelectedCard = nil
	g.firstSelectedSource = ""
	g.firstSelectedIndex = 0
	g.firstSelectedDepth = 0
}

func (g *Game) cancelSecondSelection() {
	g.secondSelectedCard = nil
	g.secondSelectedSource = ""
	g.secondSelectedIndex = 0
}

// resetPiles clears all piles and counters
func (g *Game) resetPiles() {
	g.stock = nil
	g.waste = nil
	for i := range g.foundations {
		g.foundations[i] = nil
	}
	for i := range g.tableau {
		g.tableau[i] = nil
	}
	g.moveCount = 0
	g.showDown = false
}

// toggleShow flips whether face-down cards reveal their suit and rank
func (g *Game) toggleShow() {
	g.showDown = !g.showDown
}

// shuffleAndDeal builds a 52-card deck, shuffles, deals to tableau, rest → stock
func (g *Game) shuffleAndDeal() {
	deck := make([]*Card, 0, 52)
	for s := 0; s < 4; s++ {
		for r := 1; r <= 13; r++ {
			deck = append(deck, &Card{Suit: s, Rank: r})
		}
	}
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	r.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	idx := 0
	for col := 0; col < 7; col++ {
		for row := 0; row <= col; row++ {
			c := deck[idx]
			c.FaceUp = (row == col)
			g.tableau[col] = append(g.tableau[col], c)
			idx++
		}
	}
	g.stock = deck[idx:]
}

// drawCards deals up to g.drawCount from stock → waste, or recycles waste → stock
func (g *Game) drawCards() {
	if len(g.stock) == 0 {
		for i := len(g.waste) - 1; i >= 0; i-- {
			c := g.waste[i]
			c.FaceUp = false
			g.stock = append(g.stock, c)
			g.waste = g.waste[:i]
		}
		g.moveCount++
		return
	}
	for i := 0; i < g.drawCount && len(g.stock) > 0; i++ {
		c := g.stock[len(g.stock)-1]
		g.stock = g.stock[:len(g.stock)-1]
		c.FaceUp = true
		g.waste = append(g.waste, c)
	}
	g.moveCount++ // ✅ counts as one move regardless of how many cards drawn
}
