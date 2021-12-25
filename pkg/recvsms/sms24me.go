package recvsms

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// SMS24meBackend is the backend for SMS24.me, the struct includes a name
// string
type SMS24meBackend struct {
	Name       string
	Numbers    []Number
	HTTPClient *http.Client
	// A ranking out of 10. The coder should decide this by considering
	// the backend's consistency, stability and quality. A 10 would be
	// that nearly every number works perfectly and updates the messages
	// as fast as possible, or even just actually updates the messages.
	Ranking int
}

// Returns a new backend for SMS24.me
func NewSMS24MeBackend() *SMS24meBackend {
	b := SMS24meBackend{
		Name:       "SMS24.me",
		HTTPClient: http.DefaultClient,
	}
	b.HTTPClient.Timeout = 10 * time.Second
	b.Ranking = 10
	return &b
}

func (b *SMS24meBackend) ScrapeNumbers(cache bool) ([]Number, error) {
	numbers := []Number{}
	for i := 1; i < 21; i++ {
		resp, err := http.Get("https://sms24.me/en/numbers/page/" + strconv.Itoa(i) + "/")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		str := string(bs)
		nrs := getAllStringsBetween(str, `fw-bold text-primary mb-2">`, `</div>`)
		for _, num := range nrs {
			numbers = append(numbers, Number{FullString: num})
		}
	}
	if cache {
		b.Numbers = numbers
	}
	return numbers, nil
}

func (b *SMS24meBackend) ListMessagesForNumber(n Number) ([]Message, error) {
	return nil, nil
}
