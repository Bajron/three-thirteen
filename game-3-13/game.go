package game313

import (
	"github.com/Bajron/three-thirteen/playingcards"
)

type PublicStateView struct {
	Players        []PublicPlayerView
	StartingPlayer int
	CurrentPlayer  int
	CurrentState   int
	CardsInDeck    int
	CardsOnPile    int
	PileTop        playingcards.Card
	FinalGroups    []FinalGroups
	Trumph         playingcards.Rank
}

type PublicPlayerView struct {
	CardsInHand int
	State       PlayerState
	Points      int
}

func (s State) GetPublicPlayersView() []PublicPlayerView {
	ret := make([]PublicPlayerView, len(s.Players))
	for i := 0; i < len(s.Players); i++ {
		ret[i] = PublicPlayerView{
			len(s.Players[i].Hand),
			s.Players[i].State,
			s.Players[i].Points,
		}
	}
	return ret
}

func GetPublicStateView(s *State) *PublicStateView {
	return &PublicStateView{
		s.GetPublicPlayersView(),
		s.StartingPlayer,
		s.CurrentPlayer,
		s.CurrentState,
		len(s.Deck),
		len(s.Pile),
		s.Pile.Top(),
		s.FinalGroups,
		s.Trumph,
	}
}

type MoveCommand struct {
	Player int
	Move   int
	Card   playingcards.Card
	Groups FinalGroups
}

func NewTakeCommand(p int, t int) MoveCommand {
	return MoveCommand{
		p, t, playingcards.NIL_CARD, FinalGroups{nil, nil},
	}
}

type Scores []int

type GameSession struct {
	// Current state
	state *State
	// Current trumph moves history
	history []MoveCommand
	// 3-14 history
	scoresHistory []Scores
	// All time history of 14 result
	gamesHistory []Scores
}

func NewGameSession(players int) *GameSession {
	return &GameSession{
		New(players),
		make([]MoveCommand, 0),
		make([]Scores, 0, 12),
		make([]Scores, 0),
	}
}

func (gs *GameSession) Dispatch(cmd MoveCommand) (playingcards.Card, error) {
	c := playingcards.NIL_CARD
	m := cmd.Move

	if m == TAKE_FROM_DECK || m == TAKE_FROM_PILE {
		move := Move(m)
		return gs.state.TakeMove(cmd.Player, move)
	}
	return c, &gameError{"Dispatch: not implemented"}
}
