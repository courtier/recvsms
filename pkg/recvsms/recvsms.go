package recvsms

import "strings"

type Number struct {
	CountryCode string
	PhoneNumber string
	FullString  string
}

type Message struct {
	Sender  string
	Content string
	Sent    string
}

type Backend interface {
	ScrapeNumbers() ([]Number, error)
	ListMessagesForNumber(Number) ([]Message, error)
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
