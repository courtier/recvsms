package recvsms

type SMS24meBackend struct {
	Name    string
	Numbers []Number
}

func NewSMS24MeBackend() *SMS24meBackend {
	return &SMS24meBackend{
		Name: "SMS24.me",
	}
}

func (b *SMS24meBackend) ScrapeNumbers() []Number {
	return nil
}

func (b *SMS24meBackend) ListMessagesForNumber(n Number) []string {
	return nil
}
