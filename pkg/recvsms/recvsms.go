package recvsms

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Number is a phone number listed on the backends.
// CountryCode & PhoneNumber may not always be available,
// but the FullString field will always be set.
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

// Message is a message scraped off a backend.
type Message struct {
	// Message sender if available. May not always be available.
	Sender string
	// Message content.
	Content string
	// When the message was sent. May not always be available.
	// TODO#turn this into a time.Time
	Sent string
	// When the message was scraped off/seen on the website.
	Found time.Time
}

// Backend is a SMS service that allows anyone to receive SMS represented
// as a interface. there are/will be many unique backends which are retrievable
// via the ListBackends() function.
type Backend interface {
	// ScrapeNumbers scrapes an SMS backend and returns all the numbers
	// in an array. If cache is true, the numbers will be cached and be available
	// through the GetNumbers() function.
	ScrapeNumbers(cache bool) ([]Number, error)
	// ListMessagesForNumber scrapes all the messages of the number,
	// and returns them in an array. If cache is true, messaages will
	// be cached in the Messages field.
	ListMessagesForNumber(number Number, cache bool) ([]Message, error)
	// DiffMessages() scrapes a number, then compares the newly scraped messages
	// to the cache and returns the messages that were not in the cache. Also
	// caches the new messages if cache is true.
	DiffMessagesForNumber(number Number, cache bool) ([]Message, error)
	// GetName returns the name of the backend.
	GetName() string
	// GetNumbers returns the latest cached numbers, if there are any; if not returns error.
	GetNumbers() ([]Number, error)
	// GetMessages returns the latest cached messages, if there are any; if not returns error.
	GetMessages() ([]Message, error)
	// Score is the (somewhat subjective) score of a backend, a number out of 10. The coder should decide this by considering
	// the backend's reliability, stability and quality. A 10 would be
	// that nearly every number works perfectly and updates the messages
	// as fast as possible, or even just actually updates the messages.
	Score() int
	// SetHTTPClient sets the HTTP client to be used for the backend, useful when the user
	// wants to use their own client for proxies, timeouts etc.
	SetHTTPClient(*http.Client)
}

var (
	backends = map[string]Backend{
		"SMS24.me": NewSMS24MeBackend(),
	}
)

// ListBackends lists all available backends in a map, with their names as the key.
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

func readBodyToString(body io.ReadCloser) (string, error) {
	bs, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	str := string(bs)
	return str, nil
}

func diffMessages(msgs, cache []Message) []Message {
	messages := []Message{}
	for _, m := range msgs {
		for _, c := range cache {
			if m.Sender == c.Sender && m.Content == c.Content {
				goto inCache
			}
		}
		messages = append(messages, m)
	inCache:
	}
	return messages
}
