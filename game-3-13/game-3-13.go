package game313

import (
	"github.com/Bajron/three-thirteen/playingcards"
)

type State struct {
	Players []Player
	CurrentPlayer int
	Deck playingcards.Deck
	Pile playingcards.Pile
	Trumph playingcards.Rank
}

type Player struct {
	Hand playingcards.Hand
	Points int
}

type Move int

const (
	TAKE_FROM_PILE = iota
	TAKE_FROM_DECK
)

// TODO make things private?
func New(playersNo int) (*State) {
	return &State{
		make([]Player, playersNo),
		0,
		playingcards.Create104Deck(),
		make([]playingcards.Card,0,104),
		playingcards.Rank(3)}
}

func (s *State) Deal() {
	for c := 0; c < int(s.Trumph); c++ {
		for i := 0; i < len(s.Players); i++ {
			ii := (s.CurrentPlayer + i) % len(s.Players)
			s.Players[ii].Hand =
				append(s.Players[ii].Hand, s.Deck.Draw())
		}
	}
	s.Pile = append(s.Pile, s.Deck.Draw())
}

type moveError struct {
	msg string
}

func (e *moveError) Error() string {
	return e.msg
}

func (s *State) TakeMove(player int, move Move) (playingcards.Card, error) {
	if s.CurrentPlayer != player {
		return playingcards.NIL_CARD, &moveError{"invalid move"}
	}
	switch move {
		case TAKE_FROM_PILE:
			card := s.Pile.Pop()
			s.applyTake(card)
			s.advancePlayer()
			return card, nil
		case TAKE_FROM_DECK:
			card := s.Deck.Draw()
			s.applyTake(card)
			s.advancePlayer()
			return card, nil
		default:
			return playingcards.NIL_CARD, &moveError{"invalid move"}
	}
}

func (s *State) advancePlayer() {
	s.CurrentPlayer = (s.CurrentPlayer + 1) % len(s.Players);
}

func (s *State) applyTake(card playingcards.Card) {
	s.Players[s.CurrentPlayer].Hand =
		append(s.Players[s.CurrentPlayer].Hand, card);
}

