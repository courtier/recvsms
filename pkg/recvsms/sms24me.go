package recvsms

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/biter777/countries"
)

// SMS24meBackend is the backend for SMS24.me, the struct includes a name
// string
type SMS24meBackend struct {
	Name       string
	HTTPClient *http.Client
	Numbers    []Number
	Messages   []Message
}

// NewSMS24MeBackend Returns a new backend for SMS24.me, uses a default
// HTTP client with a timeout of 10 seconds.
func NewSMS24MeBackend() *SMS24meBackend {
	b := SMS24meBackend{
		Name:       "SMS24.me",
		HTTPClient: http.DefaultClient,
	}
	b.HTTPClient.Timeout = 10 * time.Second
	return &b
}

func (b *SMS24meBackend) ScrapeNumbers(cache bool) ([]Number, error) {
	numbers := []Number{}
	for i := 1; i < 21; i++ {
		resp, err := b.HTTPClient.Get("https://sms24.me/en/numbers/page/" + strconv.Itoa(i) + "/")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		str, err := readBodyToString(resp.Body)
		if err != nil {
			return nil, err
		}
		ctrs := getAllStringsBetween(str, `<h5 class="text-secondary">`, `</h5>`)
		nrs := getAllStringsBetween(str, `fw-bold text-primary mb-2">`, `</div>`)
		for i, num := range nrs {
			country := countries.ByName(ctrs[i])
			cc := country.Info().CallCodes[0].String()
			n := num[len(cc):]
			numbers = append(numbers, Number{
				CountryCode: cc,
				PhoneNumber: n,
				FullString:  num,
			})
		}
	}
	if cache {
		b.Numbers = numbers
	}
	return numbers, nil
}

func (b *SMS24meBackend) ListMessagesForNumber(n Number, cache bool) ([]Message, error) {
	messages := []Message{}
	if len(n.FullString) < 1 {
		return nil, errors.New("invalid number")
	}
	url := "https://sms24.me/en/numbers/" + n.FullString[1:]
	resp, err := b.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	str, err := readBodyToString(resp.Body)
	if err != nil {
		return nil, err
	}
	msgs := getAllStringsBetween(str, `                    </div>

                    `, "\n            </div>")
	senders := getAllStringsBetween(str, `From: <a href="/en/messages/`, `/">`)
	for i, m := range msgs {
		messages = append(messages, Message{
			Sender:  senders[i],
			Content: m,
			Found:   time.Now(),
		})
	}
	if cache {
		b.Messages = messages
	}
	return messages, nil
}

func (b *SMS24meBackend) DiffMessagesForNumber(number Number, cache bool) ([]Message, error) {
	if b.Messages == nil {
		return nil, errors.New("empty message cache")
	}
	messages := []Message{}
	msgs, err := b.ListMessagesForNumber(number, false)
	if err != nil {
		return nil, err
	}
	for _, m := range msgs {
		for _, c := range b.Messages {
			if m.Sender == c.Sender && m.Content == c.Content {
				goto inCache
			}
		}
		messages = append(messages, m)
	inCache:
	}
	if cache {
		b.Messages = msgs
	}
	return messages, nil
}

func (b *SMS24meBackend) GetName() string {
	return b.Name
}

func (b *SMS24meBackend) GetNumbers() ([]Number, error) {
	if b.Numbers != nil {
		return b.Numbers, nil
	}
	return nil, errors.New("no cached numbers")
}

func (b *SMS24meBackend) GetMessages() ([]Message, error) {
	if b.Messages != nil {
		return b.Messages, nil
	}
	return nil, errors.New("no cached messages")
}

func (b *SMS24meBackend) Score() int {
	return 10
}

func (b *SMS24meBackend) SetHTTPClient(c *http.Client) {
	b.HTTPClient = c
}
