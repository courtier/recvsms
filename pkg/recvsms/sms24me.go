package recvsms

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/biter777/countries"
)

// SMS24meBackend is the backend for SMS24.me, the struct includes a name
// string
type SMS24meBackend struct {
	Name       string
	Numbers    []Number
	HTTPClient *http.Client
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
		ctrs := getAllStringsBetween(str, `<h5 class="text-secondary">`, `</h5>`)
		nrs := getAllStringsBetween(str, `fw-bold text-primary mb-2">`, `</div>`)
		for i, num := range nrs {
			country := countries.ByName(ctrs[i])
			cc := country.Info().CallCodes[0].String()
			n := num[len(cc):]
			numbers = append(numbers, Number{CountryCode: cc, PhoneNumber: n, FullString: num})
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

func (b *SMS24meBackend) GetName() string {
	return b.Name
}

func (b *SMS24meBackend) GetNumbers() ([]Number, error) {
	if b.Numbers != nil {
		return b.Numbers, nil
	}
	return nil, errors.New("no cached numbers")
}

func (b *SMS24meBackend) Score() int {
	return 10
}

func (b *SMS24meBackend) SetHTTPClient(c *http.Client) {
	b.HTTPClient = c
}
