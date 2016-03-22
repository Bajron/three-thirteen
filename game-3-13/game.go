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

// Player command during a game
type MoveCommand struct {
	Player int
	Move   int
	Card   playingcards.Card
	Groups FinalGroups
}

// Game marshal commands
const (
	DEAL = iota
	NEXT_TRUMPH
	NEXT_GAME
)

// Higher level game marshal commands
type GameCommand struct {
	Move int
}

func NewTakeCommand(p int, t int) MoveCommand {
	return MoveCommand{
		p, t, playingcards.NIL_CARD, FinalGroups{nil, nil},
	}
}

func NewThrowCommand(p int, c playingcards.Card) MoveCommand {
	return MoveCommand{
		p, THROW_CARD, c, FinalGroups{nil, nil},
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
	var e error
	c := playingcards.NIL_CARD
	m := cmd.Move

	e = &gameError{"Dispatch: not implemented"}

	switch m {
	case TAKE_FROM_DECK, TAKE_FROM_PILE:
		move := Move(m)
		c, e = gs.state.TakeMove(cmd.Player, move)
	case THROW_CARD:
		e = gs.state.ThrowMove(cmd.Player, cmd.Card)
	case DECLARE_DONE:
		e = gs.state.DoneMove(cmd.Player, cmd.Groups)
	case PASS_TURN:
		e = gs.state.PassMove(cmd.Player)
	}

	if e == nil {
		gs.history = append(gs.history, cmd)
	}

	return c, e
}

func (gs *GameSession) Marshal(cmd GameCommand) error {
	switch cmd.Move {
	case DEAL:
		gs.state.Deal()
		return nil
	case NEXT_TRUMPH:
		return gs.state.AdvanceRound()
	case NEXT_GAME:
	}
	return &gameError{"Marshal: not implemented"}
}

func (gs *GameSession) GetTableState() *PublicStateView {
	return GetPublicStateView(gs.state)
}

type GameServer struct {
	Sessions map[string]*GameSession
}

func NewGameServer() *GameServer {
	return &GameServer{
		make(map[string]*GameSession),
	}
}
