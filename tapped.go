package main

import "fyne.io/fyne/v2"

// Tapped handles card selection and movement
func (cw *CardWidget) Tapped(_ *fyne.PointEvent) {
	g := cw.game

	if g.firstSelectedCard == nil { // There is no first selected card

		// Prevent selecting a face-down card
		if !cw.card.FaceUp {
			return
		}
		// Prevent selecting a foundation card
		for _, pile := range g.foundations {
			if len(pile) > 0 && pile[len(pile)-1] == cw.card {
				return
			}
		}

		// Prevent selecting penultimate and double-penultimate waste cards
		if len(g.waste) > 1 && g.waste[len(g.waste)-2] == cw.card ||
			len(g.waste) > 2 && g.waste[len(g.waste)-3] == cw.card {
			return
		}
		// At this point we know that the selected card is FaceUp, but is not on a foundation
		// and is not the penultimate or double penultimate waste card.
		// So it must be either last waste card or in tableau.
		g.firstSelectedCard = cw.card

		// Find firstSelectedCard
		// Check each tableau pile
		for i, pile := range g.tableau {
			for j, c := range pile { // Checks from first to last, including FaceDown card which is no necessary
				if c == g.firstSelectedCard {
					g.firstSelectedSource = "tableau"
					g.firstSelectedIndex = i //which tableau pile
					g.firstSelectedDepth = j //which card in the pile
					g.updateUI()
					return
				}
			}
		}
		// Only remaining possibility is last waste card
		g.firstSelectedSource = "waste"
		g.firstSelectedIndex = len(g.waste) - 1
		g.firstSelectedDepth = 0
		g.updateUI()
		return

	} else { // FirstSelectedCard is already selected, meaning this is a second click
		g.secondSelectedCard = cw.card

		if g.secondSelectedCard == g.firstSelectedCard {
			// Cancel selection
			g.cancelFirstSelection()
			g.cancelSecondSelection()
			g.updateUI()
			return
		}

		// Find secondSelectedCard & make move if possible
		for i, pile := range g.tableau { // Check each tableau pile
			if len(pile) != 0 { // pile not empty
				lastCard := pile[len(pile)-1]
				if g.secondSelectedCard == lastCard && // If this line is satisfied, we have found the 2nd selected card
					g.firstSelectedCard.Rank == g.secondSelectedCard.Rank-1 &&
					g.firstSelectedCard.Color() != g.secondSelectedCard.Color() {
					if g.firstSelectedSource == "tableau" { // move from tableau to tableau
						g.tableau[i] = append(g.tableau[i], g.tableau[g.firstSelectedIndex][g.firstSelectedDepth:]...) // Move group
						g.tableau[g.firstSelectedIndex] = g.tableau[g.firstSelectedIndex][:g.firstSelectedDepth]       // Remove group from source
						// if new last card of source is face down, flip it up
						if len(g.tableau[g.firstSelectedIndex]) > 0 {
							last := g.tableau[g.firstSelectedIndex][len(g.tableau[g.firstSelectedIndex])-1]
							if !last.FaceUp {
								last.FaceUp = true
							}
						}
					} else { // move from waste
						g.tableau[i] = append(g.tableau[i], g.waste[len(g.waste)-1])
						g.waste = g.waste[:len(g.waste)-1] // Remove last card from waste
					}
					g.cancelFirstSelection()
					g.cancelSecondSelection()
					g.moveCount++
					g.updateUI()
					return
				}
			}
		}
		for i, pile := range g.foundations { // Check each foundation pile
			if len(pile) != 0 { // pile not empty
				lastCard := pile[len(pile)-1]
				if g.secondSelectedCard == lastCard &&
					g.firstSelectedCard.Rank == g.secondSelectedCard.Rank+1 &&
					g.firstSelectedCard.Suit == g.secondSelectedCard.Suit {
					if g.firstSelectedSource == "tableau" {
						g.foundations[i] = append(g.foundations[i], g.tableau[g.firstSelectedIndex][len(g.tableau[g.firstSelectedIndex])-1]) // Move card
						g.tableau[g.firstSelectedIndex] = g.tableau[g.firstSelectedIndex][:g.firstSelectedDepth]                             // Remove card from source
						if len(g.tableau[g.firstSelectedIndex]) > 0 {
							last := g.tableau[g.firstSelectedIndex][len(g.tableau[g.firstSelectedIndex])-1]
							if !last.FaceUp {
								last.FaceUp = true
							}
						}
						// g.moveCount++
					} else { // move from waste
						g.foundations[i] = append(g.foundations[i], g.waste[len(g.waste)-1])
						g.waste = g.waste[:len(g.waste)-1]
						// g.moveCount++
					}
					g.cancelFirstSelection()
					g.cancelSecondSelection()
					g.moveCount++
					g.updateUI()
					return
				}
			}
		}
	}
}

func (w *EmptyPileWidget) Tapped(_ *fyne.PointEvent) {
	g := w.game
	if g.firstSelectedCard == nil {
		return // nothing selected, do nothing
	}
	if w.pileType == "tableau" {
		// Only allow Kings to be moved to empty tableau columns
		if g.firstSelectedCard.Rank == King {
			if g.firstSelectedSource == "tableau" {
				// Move group from tableau to empty tableau
				src := g.tableau[g.firstSelectedIndex]
				group := src[g.firstSelectedDepth:]
				g.tableau[w.pileIndex] = append(g.tableau[w.pileIndex], group...)
				g.tableau[g.firstSelectedIndex] = src[:g.firstSelectedDepth]
				if len(g.tableau[g.firstSelectedIndex]) > 0 {
					last := g.tableau[g.firstSelectedIndex][len(g.tableau[g.firstSelectedIndex])-1]
					if !last.FaceUp {
						last.FaceUp = true
					}
				}
			} else if g.firstSelectedSource == "waste" {
				// Move King from waste to empty tableau
				g.tableau[w.pileIndex] = append(g.tableau[w.pileIndex], g.waste[len(g.waste)-1])
				g.waste = g.waste[:len(g.waste)-1]
			}
			g.cancelFirstSelection()
			g.cancelSecondSelection()
			g.moveCount++
			g.updateUI()
		}
	}
	if w.pileType == "foundation" {
		// Only allow Aces to be moved to empty foundation
		if g.firstSelectedCard.Rank == Ace {
			if g.firstSelectedSource == "tableau" {
				g.foundations[w.pileIndex] = append(g.foundations[w.pileIndex], g.firstSelectedCard)
				g.tableau[g.firstSelectedIndex] = g.tableau[g.firstSelectedIndex][:g.firstSelectedDepth]
				if len(g.tableau[g.firstSelectedIndex]) > 0 {
					last := g.tableau[g.firstSelectedIndex][len(g.tableau[g.firstSelectedIndex])-1]
					if !last.FaceUp {
						last.FaceUp = true
					}
				}
			} else if g.firstSelectedSource == "waste" {
				g.foundations[w.pileIndex] = append(g.foundations[w.pileIndex], g.waste[len(g.waste)-1])
				g.waste = g.waste[:len(g.waste)-1]
			}
			g.cancelFirstSelection()
			g.cancelSecondSelection()
			g.moveCount++
			g.updateUI()
		}
	}
}
