package recvsms

import (
	"io/ioutil"
	"net/http"
	"strconv"
)

type SMS24meBackend struct {
	Name    string
	Numbers []Number
}

func NewSMS24MeBackend() *SMS24meBackend {
	return &SMS24meBackend{
		Name: "SMS24.me",
	}
}

func (b *SMS24meBackend) ScrapeNumbers() ([]Number, error) {
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
	return numbers, nil
}

func (b *SMS24meBackend) ListMessagesForNumber(n Number) ([]Message, error) {
	return nil, nil
}
