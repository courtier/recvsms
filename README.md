# recvsms - use free sms services through a desktop GUI

[![Go Reference](https://pkg.go.dev/badge/github.com/courtier/recvsms.svg)](https://pkg.go.dev/github.com/courtier/recvsms/pkg/recvsms) [![Go Report Card](https://goreportcard.com/badge/github.com/courtier/recvsms)](https://goreportcard.com/report/github.com/courtier/recvsms)

You can think of recvsms as yt-dl for free sms services. It will support many SMS "backends" eventually. It also doubles as a sms receiving library.

![Backends](screenshots/backends.png)
![Numbers](screenshots/numbers.png)
![Messages](screenshots/messages.png)

## Goals:
- Support as many services as possible
- Clean implementation and library
- Easy to use GUI
- Lightweight

## TODO:
- Add more backends
- Make an actually working program
- Separate recvsms library from the actual program
- (If possible) Have a way to detect VoIP numbers
- Fix memory leak that happens when you change numbers really fast (have the same page on the right always, just update the content accordingly)