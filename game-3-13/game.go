package game313

import (
	"github.com/Bajron/three-thirteen/playingcards"
)

var GAME_COMMANDS map[int]string = nil
var PLAYER_COMMANDS map[int]string = nil
var GAME_STATES map[int]string = nil
var PLAYER_STATES map[int]string = nil

func init() {
	GAME_COMMANDS = make(map[int]string)
	GAME_COMMANDS[DEAL] = "DEAL"
	GAME_COMMANDS[NEXT_TRUMPH] = "NEXT_TRUMPH"
	GAME_COMMANDS[NEXT_GAME] = "NEXT_GAME"

	PLAYER_COMMANDS = make(map[int]string)
	PLAYER_COMMANDS[TAKE_FROM_PILE] = "TAKE_FROM_PILE"
	PLAYER_COMMANDS[TAKE_FROM_DECK] = "TAKE_FROM_DECK"
	PLAYER_COMMANDS[THROW_CARD] = "THROW_CARD"
	PLAYER_COMMANDS[PASS_TURN] = "PASS_TURN"
	PLAYER_COMMANDS[DECLARE_DONE] = "DECLARE_DONE"

	GAME_STATES = make(map[int]string)
	GAME_STATES[NOT_DEALT] = "NOT_DEALT"
	GAME_STATES[FINISHED] = "FINISHED"
	GAME_STATES[STARTING] = "STARTING"
	GAME_STATES[PLAYING] = "PLAYING"
	GAME_STATES[FINISHING] = "FINISHING"

	PLAYER_STATES = make(map[int]string)
	PLAYER_STATES[WAIT] = "WAIT"
	PLAYER_STATES[TAKE] = "TAKE"
	PLAYER_STATES[THROW] = "THROW"
	PLAYER_STATES[DONE] = "DONE"
}

type PublicStateView struct {
	Players        []PublicPlayerView
	StartingPlayer int
	DealingPlayer  int
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
		s.DealingPlayer(),
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
	Player int
	Move   int
}

func NewGameCommand(p int, m int) GameCommand {
	return GameCommand{p, m}
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

func NewPassCommand(p int) MoveCommand {
	return MoveCommand{
		p, PASS_TURN, playingcards.NIL_CARD, FinalGroups{nil, nil},
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
	// Index of priviliged player
	admin int
}

func NewGameSession(players int, admin int) *GameSession {
	return &GameSession{
		New(players),
		make([]MoveCommand, 0),
		make([]Scores, 0, 12),
		make([]Scores, 0),
		admin,
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

func (s *State) PreviousPlayer(p int) int {
	l := len(s.Players)
	return (p - 1 + l) % l
}

func (s *State) DealingPlayer() int {
	return s.PreviousPlayer(s.StartingPlayer)
}

func (gs *GameSession) Marshal(cmd GameCommand) error {
	switch cmd.Move {
	case DEAL:
		if gs.state.DealingPlayer() != cmd.Player {
			return &gameError{"Only dealer can deal"}
		}
		gs.state.Deal()
		return nil
	case NEXT_TRUMPH:
		if gs.admin != cmd.Player {
			return &gameError{"Only game admin can do this"}
		}
		return gs.state.AdvanceRound()
	case NEXT_GAME:
		if gs.admin != cmd.Player {
			return &gameError{"Only game admin can do this"}
		}
	}
	return &gameError{"Marshal: not implemented"}
}

func (gs *GameSession) GetTableState() *PublicStateView {
	return GetPublicStateView(gs.state)
}

func (gs *GameSession) GetPlayersHand(i int) playingcards.Group {
	return gs.state.Players[i].Hand
}
