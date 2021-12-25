package recvsms

type Number struct {
	CountryCode string
	PhoneNumber string
}

type Message struct {
	Sender string
}

type Backend interface {
	ScrapeNumbers() []Number
	ListMessagesForNumber(Number) []string
}

var (
	backends = map[string]Backend{
		"SMS24.me": NewSMS24MeBackend(),
	}
)

func ListBackends() map[string]Backend {
	return backends
}
