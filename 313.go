package main

import (
	"encoding/json"
	"fmt"
	"github.com/Bajron/three-thirteen/game-3-13"
	"github.com/Bajron/three-thirteen/playingcards"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SessionList struct {
	List []string
}

func jsonOrError(w http.ResponseWriter, m *Message) {
	b, err := json.Marshal(m)
	if err != nil {
		b, err = json.Marshal(NewErrorMessage(
			fmt.Sprintf("JSON encoding returned error: %v\n", err),
			m.Type,
			m.Token,
		))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	w.Write(b)
}

type GameServer struct {
	Sessions map[string]*ServerSession
}

type ServerSession struct {
	Name    string
	Session *game313.GameSession
	Players []string
}

func NewGameServer() *GameServer {
	return &GameServer{
		make(map[string]*ServerSession),
	}
}

type CodesTranslations struct {
	Cards          map[string]string
	GameStates     map[string]string
	PlayerStates   map[string]string
	GameCommands   map[string]string
	PlayerCommands map[string]string
	Suits          map[string]string
	Ranks          map[string]string
}

type SessionState struct {
	Game        *game313.PublicStateView
	Players     []string
	MovesTotal  int
	MovesOffset int
	Moves       []game313.MoveCommand
}

type SessionScores struct {
	Trumphs []game313.Scores
	Games   []game313.Scores
}

type Message struct {
	Status int
	Info   string
	Type   string
	Token  string
	Data   interface{}
}

const (
	STATUS_OK = iota
	GENERAL_ERROR
)

func NewErrorMessage(message string, inReply string, token string) *Message {
	return &Message{
		GENERAL_ERROR,
		message,
		inReply,
		token,
		nil,
	}
}

func NewOkMessage(d interface{}, inReply string, token string) *Message {
	return &Message{
		STATUS_OK,
		"OK",
		inReply,
		token,
		d,
	}
}

func (s *ServerSession) GetState() *SessionState {
	return s.GetStateWithMoves(-1)
}

func (s *ServerSession) GetStateWithMoves(since int) *SessionState {
	moves := s.Session.GetMovesHistory()
	r := &SessionState{
		s.Session.GetTableState(),
		s.Players,
		len(moves),
		len(moves),
		moves[len(moves):],
	}

	if 0 <= since && since < len(moves) {
		r.MovesOffset = since
		r.Moves = moves[since:]
	}
	return r
}

func (s *ServerSession) GetScores() *SessionScores {
	return &SessionScores{
		s.Session.GetScoresHistory(),
		s.Session.GetGamesHistory(),
	}
}

func stringifyMap(m map[int]string) map[string]string {
	r := make(map[string]string)
	for k, v := range m {
		r[fmt.Sprint(k)] = v
	}
	return r
}

func GetCodesTranslation() CodesTranslations {
	m := make(map[string]string)
	for _, s := range playingcards.SUITS {
		for _, r := range playingcards.RANKS {
			c := playingcards.Card{r, s}
			m[c.String()] = c.UString()
		}
	}

	suits := make(map[string]string)
	for _, s := range playingcards.SUITS {
		suits[fmt.Sprint(int(s))] = s.String()
	}

	ranks := make(map[string]string)
	for _, r := range playingcards.RANKS {
		ranks[fmt.Sprint(int(r))] = r.String()
	}

	return CodesTranslations{
		m,
		stringifyMap(game313.GAME_STATES),
		stringifyMap(game313.PLAYER_STATES),
		stringifyMap(game313.GAME_COMMANDS),
		stringifyMap(game313.PLAYER_COMMANDS),
		suits,
		ranks,
	}
}

func main() {
	visualTest()

	fmt.Printf("Testing test server follows, Ctrl+C to stop\n")

	gameServer := NewGameServer()

	serveMux := mux.NewRouter()
	ttRouter := serveMux.PathPrefix("/3-13/").Subrouter()
	ttRouter.StrictSlash(true)

	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Ups. Nothing special here :C\n")
	})

	serveMux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	ttRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if _, ok := q["help"]; ok {
			fmt.Fprintln(w, "<html><body><pre>")
			fmt.Fprintln(w, "Available commands: help, create, list")
			fmt.Fprintln(w, "</pre></body></html>")
		} else if values, ok := q["create"]; ok && len(values) > 0 && len(values[0]) > 0 {
			name := values[0]
			_, has := gameServer.Sessions[name]
			if has {
				fmt.Fprintf(w, "Session named like this already exists\n")
				return
			}
			pV, pOk := q["players"]
			if !pOk {
				fmt.Fprintf(w, "You need to provide players\n")
				return
			}
			players := strings.Split(pV[0], ",")

			gameServer.Sessions[name] = &ServerSession{
				name, game313.NewGameSession(len(players), 0), players,
			}

			pdata := gameServer.Sessions[name].GetState()
			jsonOrError(w, NewOkMessage(pdata, "create", ""))
		} else if _, ok := q["list"]; ok {
			l := &SessionList{make([]string, len(gameServer.Sessions))}
			i := 0
			for k := range gameServer.Sessions {
				l.List[i] = k
				i++
			}
			jsonOrError(w, NewOkMessage(l, "list", ""))
		} else if _, ok := q["translations"]; ok {
			jsonOrError(w, NewOkMessage(GetCodesTranslation(), "translations", ""))
		} else {
			fmt.Fprintln(w, "<html><body><ul>")
			for k := range gameServer.Sessions {
				fmt.Fprintf(w, "<ul><a href=\"%s\">%s</a></ul>\n", k, k)
			}
			fmt.Fprintln(w, "</ul></body></html>")
		}
	})

	ttRouter.HandleFunc("/{session}/", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["session"]
		s, ok := gameServer.Sessions[name]
		if !ok {
			fmt.Fprintf(w, "Game session does not exit\n")
			return
		}
		q := r.URL.Query()
		if _, ok := q["scores"]; ok {
			pdata := s.GetScores()
			jsonOrError(w, NewOkMessage(pdata, "session scores", ""))
		} else if v, ok := q["since"]; ok && len(v) > 0 && len(v[0]) > 0 {
			since, err := strconv.Atoi(v[0])
			if err != nil {
				msg := fmt.Sprintf("Error: %v", err)
				jsonOrError(w, NewErrorMessage(msg, "move", ""))
				return
			}
			pdata := s.GetStateWithMoves(since)
			jsonOrError(w, NewOkMessage(pdata, "session info with moves", ""))
		} else {
			pdata := s.GetState()
			jsonOrError(w, NewOkMessage(pdata, "session info", ""))
		}
	})

	ttRouter.HandleFunc("/{session}/{player}/", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["session"]
		s, ok := gameServer.Sessions[name]
		if !ok {
			fmt.Fprintf(w, "Game session does not exit\n")
			return
		}
		player := vars["player"]
		var playerIndex int
		for i, v := range s.Players {
			if v == player {
				playerIndex = i
				break
			}
		}

		q := r.URL.Query()
		if _, ok := q["help"]; ok {
			fmt.Fprintf(w, "Available commands:\n")
			fmt.Fprintf(w, "?marshal=[%d:deal|%d:nextround|%d:nextgame|]\n",
				game313.DEAL, game313.NEXT_TRUMPH, game313.NEXT_GAME)
		} else if values, ok := q["marshal"]; ok && len(values) > 0 && len(values[0]) > 0 {
			m, err := strconv.Atoi(values[0])
			if err != nil {
				msg := fmt.Sprintf("Error: %v", err)
				jsonOrError(w, NewErrorMessage(msg, "move", ""))
				return
			}

			s.Session.Marshal(game313.NewGameCommand(playerIndex, m))
			jsonOrError(w, NewOkMessage(nil, "marshal", ""))
		} else if _, ok := q["hand"]; ok {
			jsonOrError(w, NewOkMessage(s.Session.GetPlayersHand(playerIndex), "hand", ""))
		} else if values, ok := q["move"]; ok && len(values) > 0 && len(values[0]) > 0 {
			sendError := func(e error) {
				msg := fmt.Sprintf("Error: %v", e)
				jsonOrError(w, NewErrorMessage(msg, "move", ""))
			}
			m, err := strconv.Atoi(values[0])
			if err != nil {
				sendError(err)
				return
			}
			switch m {
			case game313.TAKE_FROM_DECK, game313.TAKE_FROM_PILE:
				var c playingcards.Card
				c, err = s.Session.Dispatch(game313.NewTakeCommand(playerIndex, m))
				if err != nil {
					sendError(err)
					return
				}
				jsonOrError(w, NewOkMessage(c, "move", ""))
			case game313.THROW_CARD:
				var c playingcards.Card
				if values, ok = q["card"]; ok && len(values) > 0 && len(values[0]) > 0 {
					c = playingcards.NewFromString(values[0])
				}

				if c == playingcards.NIL_CARD {
					jsonOrError(w, NewErrorMessage("invalid card parameter", "move", ""))
					return
				}

				_, err = s.Session.Dispatch(game313.NewThrowCommand(playerIndex, c))
				if err != nil {
					sendError(err)
					return
				}
				jsonOrError(w, NewOkMessage(nil, "move", ""))
			case game313.PASS_TURN:
				_, err = s.Session.Dispatch(game313.NewPassCommand(playerIndex))
				if err != nil {
					sendError(err)
					return
				}
				nextMove := len(s.Session.GetMovesHistory())
				jsonOrError(w, NewOkMessage(nextMove, "move", ""))
			case game313.DECLARE_DONE:
				fg := game313.NewFinalGroups()
				if values, ok = q["groups"]; ok && len(values) > 0 && len(values[0]) > 0 {
					aGroups := strings.Split(values[0], ";")
					for _, g := range aGroups {
						aGroup := strings.Split(g, ",")
						group := make(playingcards.Group, len(aGroup))
						for i, aC := range aGroup {
							group[i] = playingcards.NewFromString(aC)
						}
						fg.Set = append(fg.Set, group)
					}
				}
				if values, ok = q["unassigned"]; ok && len(values) > 0 && len(values[0]) > 0 {
					aGroup := strings.Split(values[0], ",")
					for _, aC := range aGroup {
						fg.Unassigned = append(fg.Unassigned, playingcards.NewFromString(aC))
					}
				}
				_, err = s.Session.Dispatch(game313.NewDoneCommand(playerIndex, fg))
				if err != nil {
					sendError(err)
					return
				}
				nextMove := len(s.Session.GetMovesHistory())
				jsonOrError(w, NewOkMessage(nextMove, "move", ""))
			}

		} else {
			pdata := s.GetState()
			jsonOrError(w, NewOkMessage(pdata, "FIXME", ""))
		}
	})

	server := &http.Server{
		Addr:           ":8080",
		Handler:        serveMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := server.ListenAndServeTLS("cert.pem", "key.pem")
	fmt.Printf("Bad error: %v\n", err)
}

func visualTest() {
	tmp := playingcards.Card{2, 1}
	fmt.Printf("Hello, world.\n")
	fmt.Printf("Test card: %s.\n", tmp)
	fmt.Printf("Test card: %s.\n", tmp.UString())
	fmt.Printf("Test card: %s.\n", tmp.StringWithUnicode())

	fmt.Printf("Clubs: %d\n", playingcards.CLUBS)

	for _, s := range playingcards.SUITS {
		c := playingcards.Card{playingcards.JACK, s}
		fmt.Printf(" %s %s,", c.UString(), c.StringWithUnicode())
	}
	fmt.Printf("\n")

	fmt.Printf("Take a look at the cards:\n")
	for _, s := range playingcards.SUITS {
		for _, r := range playingcards.RANKS {
			fmt.Printf(" %s", playingcards.Card{r, s}.UString())
		}
	}
	fmt.Printf("\nAll ok?\n")
}
