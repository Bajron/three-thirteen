package playingcards

import "math/rand"

type Deck []Card
type Pile []Card

func Create52Deck() Deck {
	ret := make([]Card, 52)
	i := 0
	for _, suit := range SUITS {
		for _, rank := range RANKS[1:] {
			ret[i] = Card{rank, suit}
			i++
		}
	}
	return ret
}

func Create104Deck() Deck {
	ret := make([]Card, 104)
	d52 := Create52Deck()
	copy(ret, d52)
	copy(ret[52:], d52)
	return ret
}

func findCard(cards []Card, c Card) int {
	for i, v := range cards {
		if v == c {
			return i
		}
	}
	return -1
}

func (deck Deck) Find(c Card) int {
	return findCard(deck, c)
}

func (deck *Deck) Draw() (c Card) {
	if len(*deck) <= 0 {
		return NIL_CARD
	}
	c = (*deck)[0]
	*deck = (*deck)[1:]
	return
}

func (deck Deck) Shuffle() {
	for i := len(deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
}

func (pile *Pile) Push(card Card) {
	*pile = append(*pile, card)
}

func (pile *Pile) Pop() Card {
	c, idx := pile.getLast()
	if idx >= 0 {
		*pile = (*pile)[:idx]
	}
	return c
}

func (pile Pile) Top() Card {
	c, _ := pile.getLast()
	return c
}

func (pile Pile) getLast() (c Card, i int) {
	if len(pile) <= 0 {
		return NIL_CARD, -1
	}
	idx := len(pile) - 1
	c = (pile)[idx]
	return c, idx
}
