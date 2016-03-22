package game313

import (
	"fmt"
	"github.com/Bajron/three-thirteen/playingcards"
	"math/rand"
)

// Game states
const (
	NOT_DEALT = iota
	FINISHED
	STARTING
	PLAYING
	FINISHING
)

// Player states
const (
	WAIT = iota
	TAKE
	THROW
	DONE
)

type State struct {
	Players        []Player
	StartingPlayer int
	CurrentPlayer  int
	CurrentState   int
	Deck           playingcards.Deck
	Pile           playingcards.Pile
	FinalGroups    []FinalGroups
	Trumph         playingcards.Rank
}

type Player struct {
	Hand   playingcards.Group
	State  PlayerState
	Points int
}

type FinalGroups struct {
	Set        []playingcards.Group
	Unassigned playingcards.Group
}

type PlayerState int
type Move int

// Move types
const (
	TAKE_FROM_PILE = iota
	TAKE_FROM_DECK
	THROW_CARD
	PASS_TURN
	DECLARE_DONE
)

// TODO make things private?
func New(playersNo int) *State {
	var startingPlayer int = rand.Int() % playersNo
	return newWithStaringPlayer(playersNo, startingPlayer)
}

func newWithStaringPlayer(playersNo int, startingPlayer int) *State {
	ret := &State{
		make([]Player, playersNo),
		startingPlayer, startingPlayer,
		NOT_DEALT,
		playingcards.Create104Deck(),
		make([]playingcards.Card, 0, 104),
		make([]FinalGroups, playersNo),
		playingcards.Rank(3)}
	ret.Players[startingPlayer].State = TAKE
	return ret
}

func (s *State) CardsToDeal() int {
	if s.Trumph != playingcards.ACE {
		return int(s.Trumph)
	}
	return 14
}

func (s *State) Deal() {
	if s.CurrentState != NOT_DEALT {
		return
	}
	ctd := s.CardsToDeal()
	for c := 0; c < ctd; c++ {
		for i := 0; i < len(s.Players); i++ {
			ii := (s.CurrentPlayer + i) % len(s.Players)
			s.Players[ii].Hand =
				append(s.Players[ii].Hand, s.Deck.Draw())
		}
	}
	s.Pile = append(s.Pile, s.Deck.Draw())
	s.CurrentState = STARTING
}

type moveError struct {
	msg string
}

func (e *moveError) Error() string {
	return e.msg
}

type gameError struct {
	msg string
}

func (e *gameError) Error() string {
	return e.msg
}

func (s *State) getJockerMatch() playingcards.RankMatch {
	return func(r playingcards.Rank) bool {
		return s.Trumph == r
	}
}

func (s *State) TakeMove(player int, move Move) (playingcards.Card, error) {
	if s.CurrentState == FINISHED {
		return playingcards.NIL_CARD, &moveError{"you cannot take, round is finished"}
	}
	if s.CurrentPlayer != player {
		return playingcards.NIL_CARD, &moveError{"wrong player"}
	}
	if s.Players[player].State != TAKE {
		return playingcards.NIL_CARD, &moveError{"already took"}
	}
	if int(s.Trumph) < len(*s.currentPlayerHand()) {
		return playingcards.NIL_CARD, &moveError{"already took"}
	}

	switch move {
	case TAKE_FROM_PILE:
		card := s.Pile.Pop()
		s.applyTake(card)
		return card, nil
	case TAKE_FROM_DECK:
		card := s.Deck.Draw()
		s.applyTake(card)
		return card, nil
	default:
		return playingcards.NIL_CARD, &moveError{"invalid move"}
	}
}

func (s *State) ThrowMove(player int, card playingcards.Card) error {
	if s.CurrentPlayer != player {
		return &moveError{"invalid player"}
	}
	if !s.currentPlayerHand().Has(card) {
		return &moveError{"throwing not owned card"}
	}
	if int(s.Trumph) >= len(*s.currentPlayerHand()) {
		return &moveError{"not enough cards to throw"}
	}

	s.currentPlayerHand().Extract(card)
	s.Pile.Push(card)
	s.Players[player].State = DONE

	if len(s.Deck) == 0 {
		s.shufflePileForDeck()
	}

	return nil
}

func (s *State) shufflePileForDeck() {
	s.Deck = append(s.Deck, s.Pile[:len(s.Pile)-1]...)
	s.Deck.Shuffle()
}

func (s *State) PassMove(player int) error {
	if s.CurrentState == FINISHING {
		return &moveError{"you cannot pass during final round"}
	}
	if s.CurrentPlayer != player {
		return &moveError{"move from invalid player"}
	}
	if s.Players[player].State != DONE {
		return &moveError{"you need to do some moves first"}
	}

	s.advancePlayer()
	return nil
}

func (s *State) DoneMove(player int, groups FinalGroups) error {
	if s.CurrentState == STARTING {
		return &moveError{"you cannot finish in first round"}
	}
	if s.CurrentPlayer != player {
		return &moveError{"it is not your turn"}
	}
	if s.Players[player].State != DONE {
		return &moveError{"it should be last action in turn"}
	}
	if len(groups.Unassigned) > 0 && s.CurrentState != FINISHING {
		return &moveError{"first done cannot have unassigned cards"}
	}
	received := make(playingcards.Group, 0, int(s.Trumph))
	for _, g := range groups.Set {
		if !playingcards.IsSetOrSeq(g, s.getJockerMatch()) {
			return &moveError{"invalid set provided"}
		}
		received = append(received, g...)
	}
	received = append(received, groups.Unassigned...)
	if ll, lr := len(received), len(s.Players[player].Hand); ll != lr {
		return &moveError{
			fmt.Sprintf("card number does not match %d!=%d", ll, lr)}
	}
	for _, c := range received {
		if !s.Players[player].Hand.Has(c) {
			return &moveError{
				fmt.Sprintf("not owned card detected %v", c)}
		}
	}

	s.FinalGroups[player] = groups
	s.Players[player].Hand = s.Players[player].Hand[:0]
	s.CurrentState = FINISHING
	s.advancePlayer()

	return nil
}

func (s *State) AdvanceRound() error {
	if s.CurrentState != FINISHED {
		return &gameError{"round is not finished yet"}
	}
	for i := 0; i < len(s.Players); i++ {
		for _, c := range s.FinalGroups[i].Unassigned {
			if c.Rank > 10 {
				s.Players[i].Points += 10
			} else {
				s.Players[i].Points += int(c.Rank)
			}
		}
	}
	s.FinalGroups = make([]FinalGroups, len(s.Players))
	if s.Trumph != playingcards.ACE {
		if s.Trumph < playingcards.KING {
			s.Trumph++
		} else {
			s.Trumph = playingcards.ACE
		}
		s.Deck = playingcards.Create104Deck()
	} else {
		s.Deck = make(playingcards.Deck, 0)
	}
	s.Pile = s.Pile[:0]
	s.CurrentState = NOT_DEALT
	s.StartingPlayer = (s.StartingPlayer + 1) % len(s.Players)
	return nil
}

func (s *State) currentPlayerHand() *playingcards.Group {
	return &s.Players[s.CurrentPlayer].Hand
}

func (s *State) advancePlayer() {
	s.Players[s.CurrentPlayer].State = WAIT
	s.CurrentPlayer = (s.CurrentPlayer + 1) % len(s.Players)
	s.Players[s.CurrentPlayer].State = TAKE
	if s.CurrentState == STARTING &&
		s.CurrentPlayer == s.StartingPlayer {
		s.CurrentState = PLAYING
	}
	if s.CurrentState == FINISHING &&
		s.CurrentPlayer == s.StartingPlayer {
		s.CurrentState = FINISHED
	}
}

func (s *State) applyTake(card playingcards.Card) {
	s.Players[s.CurrentPlayer].Hand =
		append(s.Players[s.CurrentPlayer].Hand, card)
	s.Players[s.CurrentPlayer].State = THROW
}
