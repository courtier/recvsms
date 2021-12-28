package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/courtier/recvsms/pkg/recvsms"
)

type topMinBottomMax struct{}
type extendBottomList struct{}

var (
	a                fyne.App
	w                fyne.Window
	scrapeInProgress = sync.Mutex{}
)

func main() {
	a = app.NewWithID("com.github.courtier.recvsms")

	w = a.NewWindow("recvsms")
	w.Resize(fyne.NewSize(640, 480))

	rand.Seed(time.Now().UnixMicro())

	backendPickerScreen()

	w.ShowAndRun()
}

func backendPickerScreen() {
	backends := recvsms.ListBackends()
	backendNames := recvsms.BackendNames()
	backendLength := recvsms.BackendsLength()

	prompt := widget.NewLabel("Pick how you would like the numbers to be scraped:")
	prompt.TextStyle.Bold = true

	info := widget.NewLabel(fmt.Sprintf("%d backend(s) available.", backendLength))
	info.TextStyle.Bold = true

	progress := widget.NewProgressBar()
	progress.Min = 0
	progress.Max = float64(len(backendNames))
	progress.Hide()

	output := widget.NewTextGrid()
	output.Hide()

	optionNames := make([]string, len(backendNames))
	for i, b := range backendNames {
		optionNames[i] = b + fmt.Sprintf(" (Score: %d)", backends[b].Score())
	}
	backendPicker := widget.NewCheckGroup(optionNames, func(picked []string) {})

	backendPickerScroll := container.NewVScroll(backendPicker)
	backendPickerScroll.SetMinSize(backendPicker.MinSize())

	pickAll := widget.NewCheck("All", func(picked bool) {
		if picked {
			backendPicker.SetSelected(backendPicker.Options)
		} else {
			backendPicker.SetSelected(nil)
		}
	})

	errorText := widget.NewTextGrid()

	errorAlert := widget.NewPopUp(errorText, w.Canvas())
	errorAlert.Hide()

	startScrape := widget.NewButton("Start Scraping", func() {
		if len(backendPicker.Selected) < 1 {
			errorText.SetText("No backends selected")
			errorAlert.Move(alignMiddle(errorAlert.MinSize().Width, errorAlert.MinSize().Height))
			errorAlert.Show()
			return
		}
		numbers := scrapeNumbers(backendPicker.Selected, progress, output)
		if len(numbers) < 1 {
			fmt.Fprintln(os.Stderr, "no numbers could be scraped")
			output.SetText(output.Text() + "\nError: No numbers could be scraped")
			return
		}
		// TODO: perhaps add this as an option so if the user started multiple scrapes, we won't override numbers afterwards
		listNumbersScreen(numbers)
	})

	w.SetContent(container.NewVBox(
		info,
		prompt,
		pickAll,
		backendPickerScroll,
		startScrape,
		progress,
		output,
	))
}

func listNumbersScreen(numbers []*recvsms.Number) {
	numberPageRight := container.New(&topMinBottomMax{})
	numberPageRight.Hide()

	numberListLeft := widget.NewList(
		func() int {
			return len(numbers)
		},
		func() fyne.CanvasObject {
			card := widget.NewCard("Blank Number", "Blank Country - Blank Backend", nil)
			card.Resize(card.MinSize())
			return card
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Card).SetTitle(numbers[i].FullString)
			o.(*widget.Card).SetSubTitle(numbers[i].CountryCode + " - " + numbers[i].Backend.GetName())
		},
	)
	numberListLeft.OnSelected = func(i widget.ListItemID) {
		displayNumberPage(numbers[i], numberPageRight)
	}

	// TODO: add a way to hide/show certain backends/countries
	// TODO: add search
	sortByListOptions := []string{"Country", "Backend", "Backend Score"}
	sortByList := widget.NewSelect(sortByListOptions, func(s string) {
		fmt.Println(s)
	})
	sortByList.PlaceHolder = "Sort By"
	pickRandomButton := widget.NewButton("Random Number", func() {
		numberListLeft.Select(rand.Intn(numberListLeft.Length()))
	})
	horiButtons := container.NewHBox(pickRandomButton, sortByList)

	leftSide := container.New(&topMinBottomMax{}, horiButtons, numberListLeft)

	split := container.NewHSplit(leftSide, numberPageRight)
	split.SetOffset(0)

	w.SetContent(split)
}

func displayNumberPage(n *recvsms.Number, c *fyne.Container) {
	messageList := container.New(&extendBottomList{})
	messageList.Hide()
	prepareMessageListForNumber(n, messageList)

	copyButton := widget.NewButton("Copy Number", func() {
		w.Clipboard().SetContent(n.FullString)
	})
	refreshButton := widget.NewButton("Refresh Messages", func() {
		createMessagesList(n, messageList)
	})
	numberButtons := container.NewCenter(container.NewHBox(copyButton, refreshButton))
	numberLabel := widget.NewLabel(n.FullString + " - " + getNumbersCountry(n) + " - " + n.Backend.GetName())
	numberLabel.Alignment = fyne.TextAlignCenter

	numberInfoTop := container.NewVBox(container.NewBorder(numberLabel, numberButtons, nil, nil))

	c.Objects = []fyne.CanvasObject{numberInfoTop, messageList}
	c.Show()
}

func prepareMessageListForNumber(n *recvsms.Number, c *fyne.Container) {
	if n.Messages != nil && len(n.Messages) > 0 {
		c.Objects = []fyne.CanvasObject{messageListToList(n.Messages)}
	} else {
		c.Objects = []fyne.CanvasObject{container.NewCenter(widget.NewButton("Scrape Messages", func() {
			createMessagesList(n, c)
		}))}
	}
	c.Show()
}

func createMessagesList(n *recvsms.Number, c *fyne.Container) {
	messages, err := n.Backend.ListMessagesForNumber(n, true)
	if err != nil {
		c.Objects = []fyne.CanvasObject{widget.NewTextGridFromString(err.Error())}
	} else if len(messages) > 0 {
		c.Objects = []fyne.CanvasObject{messageListToList(messages)}
	} else {
		c.Objects = []fyne.CanvasObject{widget.NewTextGridFromString("No messages found.")}
	}
}

func messageListToList(messages []*recvsms.Message) *widget.List {
	return widget.NewList(
		func() int {
			return len(messages)
		},
		func() fyne.CanvasObject {
			card := widget.NewCard("Sender - Sent", "Message Content", nil)
			card.Resize(card.MinSize())
			return card
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			var s string
			if messages[i].Sent == "" {
				s = "Sent: N/A"
			} else {
				s = messages[i].Sent
			}
			o.(*widget.Card).SetTitle(messages[i].Sender + " - " + s)
			o.(*widget.Card).SetSubTitle(messages[i].Content)
		},
	)
}

func scrapeNumbers(backends []string, progress *widget.ProgressBar, output *widget.TextGrid) []*recvsms.Number {
	defer scrapeInProgress.Unlock()
	scrapeInProgress.Lock()
	progress.Show()
	output.Show()
	nbrChan, beChan := make(chan *recvsms.Number, len(backends)), make(chan error)
	for _, be := range backends {
		be = strings.Split(be, " (Score: ")[0]
		b := recvsms.ListBackends()[be]
		go func(backend recvsms.Backend, nbrChan chan *recvsms.Number, beChan chan error) {
			nbrs, err := backend.ScrapeNumbers(true)
			if err != nil {
				beChan <- err
				return
			}
			for _, n := range nbrs {
				nbrChan <- n
			}
			beChan <- nil
		}(b, nbrChan, beChan)
	}
	counter := 0
	numbers := []*recvsms.Number{}
receiveLoop:
	for {
		select {
		case nbr := <-nbrChan:
			numbers = append(numbers, nbr)
		case err := <-beChan:
			counter++
			progress.SetValue(progress.Value + 1)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				output.SetText(output.Text() + "\nError:" + err.Error())
			}
			if counter == len(backends) {
				break receiveLoop
			}
		}
	}
	return numbers
}

func alignMiddle(minWidth, minHeight float32) fyne.Position {
	return fyne.NewPos((w.Content().Size().Width/2)-(minWidth/2), (w.Content().Size().Height/2)-(minHeight/2))
}

func (t *topMinBottomMax) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) == 0 {
		return fyne.NewSize(0, 0)
	}
	w, h := float32(0), float32(0)
	top := objects[0]
	w += top.MinSize().Width
	h += top.MinSize().Height
	bottom := objects[1]
	h += bottom.Size().Height
	return fyne.NewSize(w, h)
}

func (t *topMinBottomMax) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) == 0 {
		return
	}
	pos := fyne.NewPos(0, 0)
	top := objects[0]
	top.Resize(fyne.NewSize(containerSize.Width, top.MinSize().Height))
	top.Move(pos)
	pos = pos.Add(fyne.NewPos(0, top.Size().Height))
	bottom := objects[1]
	bottom.Resize(fyne.NewSize(containerSize.Width, containerSize.Height-top.MinSize().Height))
	bottom.Move(pos)
}

func (e *extendBottomList) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) == 0 {
		return fyne.NewSize(0, 0)
	}
	return objects[0].Size()
}

func (e *extendBottomList) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) == 0 {
		return
	}
	pos := fyne.NewPos(0, 0)
	top := objects[0]
	top.Resize(fyne.NewSize(containerSize.Width, containerSize.Height))
	top.Move(pos)
}

func getNumbersCountry(n *recvsms.Number) string {
	if n.CountryName == "" {
		return n.CountryCode
	}
	return n.CountryName
}
