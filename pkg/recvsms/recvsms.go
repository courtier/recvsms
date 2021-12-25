package recvsms

import (
	"net/http"
	"strings"
	"time"
)

type Number struct {
	// Country code.
	CountryCode string
	// Phone number without country code.
	PhoneNumber string
	// String including the country code and the number.
	// Probably no space in-between. This should only be used when
	// we cannot separate CC and number.
	FullString string
}

type Message struct {
	// Message sender if available
	Sender string
	// Message content
	Content string
	// When the message was sent.
	// TODO#turn this into a time.Time
	Sent string
	// When the message was scraped off/seen on the website.
	Found time.Time
}

type Backend interface {
	// ScrapeNumbers scrapes an SMS backend and returns all the numbers
	// in an array. If cache is true, these will be saved in the Numbers
	// field of the backend struct, which can be accessed directly by the user.
	// Cache can be set to false, to conserve memory.
	ScrapeNumbers(cache bool) ([]Number, error)
	// ListMessagesForNumber scrapes all the messages of the number,
	// and returns them in an array. Not mature enough to write a definite
	// description. TODO#
	ListMessagesForNumber(Number) ([]Message, error)
	// Returns the name of the backend.
	GetName() string
	// Returns the latest cached numbers, if there are any; if not returns error.
	GetNumbers() ([]Number, error)
	// (Subjective) score of a backend, a number out of 10. The coder should decide this by considering
	// the backend's reliability, stability and quality. A 10 would be
	// that nearly every number works perfectly and updates the messages
	// as fast as possible, or even just actually updates the messages.
	Score() int
	// Set the HTTP client to be used for the backend, useful when the user
	// wants to use their own client for proxies, timeouts etc.
	SetHTTPClient(*http.Client)
}

var (
	backends = map[string]Backend{
		"SMS24.me": NewSMS24MeBackend(),
	}
)

func ListBackends() map[string]Backend {
	return backends
}

func getAllStringsBetween(str, left, right string) []string {
	strs := []string{}
	b := strings.Split(str, left)
	if len(b) > 1 {
		for _, s := range b[1:] {
			d := strings.Split(s, right)
			if len(d) > 0 {
				strs = append(strs, d[0])
			}
		}
	}
	return strs
}
