package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// Suits
const (
	Hearts = iota
	Spades
	Diamonds
	Clubs
)

// Ranks
const (
	Ace   = 1
	Jack  = 11
	Queen = 12
	King  = 13
)

// Visual constants
var (
	cardSize      = fyne.NewSize(120, 180)
	cardOverlap   = float32(45)
	wasteOffsetX  = float32(52)
	cardBackColor = color.NRGBA{R: 0xB0, G: 0xD0, B: 0xFF, A: 0xFF} // Light blue for card backs
)

// Card holds suit, rank, face-up state, and its widget
type Card struct {
	Suit   int
	Rank   int
	FaceUp bool
	Widget *CardWidget
}

// CardWidget is a draggable card display
type CardWidget struct {
	widget.BaseWidget // base widget for Fyne which is, by this line, embedded into CardWidget
	card              *Card
	game              *Game
}

// NewCardWidget wraps a Card into a widget
func NewCardWidget(c *Card, g *Game) *CardWidget {
	cw := &CardWidget{card: c, game: g}
	cw.ExtendBaseWidget(cw)
	c.Widget = cw
	return cw
}

type EmptyPileWidget struct {
	widget.BaseWidget
	game      *Game
	pileType  string // "tableau" or "foundation"
	pileIndex int
}

func NewEmptyPileWidget(game *Game, pileType string, pileIndex int) *EmptyPileWidget {
	w := &EmptyPileWidget{game: game, pileType: pileType, pileIndex: pileIndex}
	w.ExtendBaseWidget(w)
	return w
}

func (c *Card) Color() string {
	if c.Suit == 0 || c.Suit == 2 {
		return "red"
	} else {
		return "black"
	}
}

// cardRenderer lays out & refreshes a CardWidget
type cardRenderer struct {
	cw   *CardWidget
	rect *canvas.Rectangle
	text *canvas.Text
	objs []fyne.CanvasObject
}

func (r *cardRenderer) Layout(size fyne.Size) {
	r.rect.Resize(size)
	r.text.Move(fyne.NewPos(8, 8))
}

func (r *cardRenderer) MinSize() fyne.Size {
	return cardSize
}

func (r *cardRenderer) Refresh() {
	lbl, clr := r.cw.cardFace()
	r.text.Text = lbl
	r.text.Color = clr
	r.text.Refresh()
}

func (r *cardRenderer) Objects() []fyne.CanvasObject { return r.objs }
func (r *cardRenderer) Destroy()                     {}

// CreateRenderer draws the card face or back
func (cw *CardWidget) CreateRenderer() fyne.WidgetRenderer {
	var rect *canvas.Rectangle
	if cw.card.FaceUp {
		rect = canvas.NewRectangle(color.White)
	} else if cw.game.showDown {
		rect = canvas.NewRectangle(color.NRGBA{R: 0xB0, G: 0xE0, B: 0xFF, A: 0xFF}) // light blue for unmasked face-down
	} else {
		rect = canvas.NewRectangle(cardBackColor)
	}
	rect.SetMinSize(cardSize)
	rect.StrokeColor = color.Black
	rect.StrokeWidth = 1

	// Highlight if selected or part of selected group in tableau
	g := cw.game
	if g.firstSelectedCard != nil {
		// If this card is the selected card, always highlight
		if g.firstSelectedCard == cw.card {
			rect.StrokeColor = color.NRGBA{R: 0x00, G: 0x80, B: 0xFF, A: 0xFF} // blue
			rect.StrokeWidth = 3
		} else if g.firstSelectedSource == "tableau" {
			// Find this card's column and index
			for colIdx, pile := range g.tableau {
				for cardIdx, c := range pile {
					if c == cw.card && colIdx == g.firstSelectedIndex && cardIdx >= g.firstSelectedDepth {
						rect.StrokeColor = color.NRGBA{R: 0x00, G: 0x80, B: 0xFF, A: 0xFF} // blue
						rect.StrokeWidth = 3
					}
				}
			}
		}
	}

	label, clr := cw.cardFace()
	text := canvas.NewText(label, clr)
	text.TextSize = 24 // Adjusted font size for card face
	objs := []fyne.CanvasObject{rect, text}
	return &cardRenderer{cw: cw, rect: rect, text: text, objs: objs}
}

// cardFace returns label and color for face-up or face-down card
func (cw *CardWidget) cardFace() (string, color.Color) {
	ranks := map[int]string{1: "A", 11: "J", 12: "Q", 13: "K"}
	r := fmt.Sprint(cw.card.Rank)
	if s, ok := ranks[cw.card.Rank]; ok {
		r = s
	}
	suits := []string{"♥", "♠", "♦", "♣"}
	sym := suits[cw.card.Suit]
	lbl := r + sym

	if cw.card.FaceUp {
		if cw.card.Suit == Hearts || cw.card.Suit == Diamonds {
			return lbl, color.NRGBA{R: 0xE0, G: 0x00, B: 0x00, A: 0xFF}
		}
		return lbl, color.Black
	}

	if cw.game.showDown {
		return lbl, color.Gray{Y: 0x80}
	}
	return "", color.Gray{Y: 0x80}
}

var _ fyne.Tappable = (*CardWidget)(nil)

func (w *EmptyPileWidget) CreateRenderer() fyne.WidgetRenderer {
	rect := canvas.NewRectangle(color.NRGBA{R: 0xDD, G: 0xDD, B: 0xDD, A: 0xFF})
	rect.SetMinSize(cardSize)
	return widget.NewSimpleRenderer(rect)
}
