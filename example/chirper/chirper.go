package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var bindAddr = flag.String("http", ":9090", "HTTP bind address")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "chirper is an HTTP server that implements a simple Twitter-like API.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\n")
		fmt.Fprintf(os.Stderr, "\tchirper [options]\n\n")
		fmt.Fprintf(os.Stderr, "The options are:\n\n")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}
	flag.Parse()

	go startChirping()

	http.HandleFunc("/chirps", listChirps)

	log.Printf("Starting server on %s\n", *bindAddr)
	err := http.ListenAndServe(*bindAddr, nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %s", err)
	}
}

type chirp struct {
	Text  string
	Score int
	Date  time.Time
}

var numChirps = 5
var chirps = make([]*chirp, numChirps)
var lastModified = time.Now()

func startChirping() {
	for i := 0; true; i = (i + 1) % len(lines) {
		chirps[i%cap(chirps)] = &chirp{
			Text:  lines[i],
			Score: rand.Intn(5),
			Date:  time.Now(),
		}

		// Vote up chirps.
		for _, c := range chirps {
			if c != nil {
				c.Score += rand.Intn(3)
			}
		}

		lastModified = time.Now()
		time.Sleep(750 * time.Millisecond)
	}
}

func listChirps(w http.ResponseWriter, r *http.Request) {
	n := len(chirps)
	for i, c := range chirps {
		if c == nil {
			n = i
			break
		}
	}

	w.Header().Add("cache-control", "public, max-age=3, s-maxage=3")
	w.Header().Add("Vary", "Accept-Encoding")
	w.Header().Add("Last-Modified", lastModified.Format(http.TimeFormat))

	data, err := json.MarshalIndent(chirps[:n], "", "  ")
	if err != nil {
		log.Fatalf("JSON encoding failed: %s", err)
	}
	w.Write(data)

	log.Printf("Listed chirps")
}

// lines are from Shakespeare's sonnets XVIII and CXVI.
var lines = []string{
	"Shall I compare thee to a summer's day?",
	"Thou art more lovely and more temperate:",
	"Rough winds do shake the darling buds of May,",
	"And summer's lease hath all too short a date:",
	"Sometime too hot the eye of heaven shines,",
	"And often is his gold complexion dimm'd;",
	"And every fair from fair sometime declines,",
	"By chance or nature's changing course untrimm'd;",
	"But thy eternal summer shall not fade",
	"Nor lose possession of that fair thou owest;",
	"Nor shall Death brag thou wander'st in his shade,",
	"When in eternal lines to time thou growest:",
	"So long as men can breathe or eyes can see,",
	"So long lives this and this gives life to thee.",

	"Let me not to the marriage of true minds",
	"Admit impediments. Love is not love",
	"Which alters when it alteration finds,",
	"Or bends with the remover to remove:",
	"O no! it is an ever-fixed mark",
	"That looks on tempests and is never shaken;",
	"It is the star to every wandering bark,",
	"Whose worth's unknown, although his height be taken.",
	"Love's not Time's fool, though rosy lips and cheeks",
	"Within his bending sickle's compass come:",
	"Love alters not with his brief hours and weeks,",
	"But bears it out even to the edge of doom.",
	"If this be error and upon me proved,",
	"I never writ, nor no man ever loved.",
}
