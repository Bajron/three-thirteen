package main

import (
	"encoding/json"
	"fmt"
	"github.com/Bajron/three-thirteen/game-3-13"
	"github.com/Bajron/three-thirteen/playingcards"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
)

type SessionList struct {
	List []string
}

func jsonOrError(w http.ResponseWriter, x interface{}) {
	b, err := json.Marshal(x)
	if err == nil {
		w.Write(b)
	} else {
		fmt.Fprintf(w, "JSON encoding returned error: %v\n", err)
	}
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

	ttRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if values, ok := q["create"]; ok && len(values) > 0 && len(values[0]) > 0 {
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
				name, game313.NewGameSession(len(players)), players,
			}

			pdata := gameServer.Sessions[name].Session.GetTableState()
			jsonOrError(w, pdata)
		} else if _, ok := q["list"]; ok {
			l := &SessionList{make([]string, len(gameServer.Sessions))}
			i := 0
			for k := range gameServer.Sessions {
				l.List[i] = k
				i++
			}
			jsonOrError(w, l)
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
		pdata := s.Session.GetTableState()
		jsonOrError(w, pdata)
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
		for i, v := range s.Players {
			if v == player {
				fmt.Fprintf(w, "Player index: %d\n", i)
			}
		}
		pdata := s.Session.GetTableState()
		jsonOrError(w, pdata)
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
