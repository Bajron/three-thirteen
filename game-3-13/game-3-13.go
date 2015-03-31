package game313

import (
	"github.com/Bajron/three-thirteen/playingcards"
)

type State struct {
	Players []Player
	Deck playingcards.Deck
	Pile playingcards.Deck
	Trumph playingcards.Rank
}

type Player struct {
	Hand playingcards.Hand
	Points int
}

func New(playersNo int) (*State) {
	return &State{
		make([]Player, playersNo),
		playingcards.Create104Deck(),
		make([]playingcards.Card,0,104),
		playingcards.Rank(3)}
}

func (s *State) Deal() {
	for c := 0; c < int(s.Trumph); c++ {
		for i := range s.Players {
			s.Players[i].Hand =
				append(s.Players[i].Hand, s.Deck.Draw())
		}
	}
	s.Pile = append(s.Pile, s.Deck.Draw())
}

