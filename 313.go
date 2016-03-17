package main

import (
	"fmt"
	"github.com/Bajron/three-thirteen/playingcards"
	"net/http"
	"time"
)

func main() {
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

	fmt.Printf("Testing test server follows, Ctrl+C to stop\n")

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Ups. Nothing special here :C\n")
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
