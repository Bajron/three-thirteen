package playingcards 

import "fmt"

type Card struct {
	Rank Rank
	Suit Suit
}

type Suit int8 

const (
	CLUBS = iota
	DIAMONDS
	SPADES
	HEARTS
)

type Rank int8 

const (
	ACE = iota
	JACK = 11
	QUEEN
	KING
)

func (c Card) String() string {
	return fmt.Sprintf("%d %d", c.Rank, c.Suit)
}
