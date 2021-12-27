package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/courtier/recvsms/pkg/recvsms"
)

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

	backendPicker := widget.NewCheckGroup(backendNames, func(picked []string) {})

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

func listNumbersScreen(numbers []recvsms.Number) {
	numberPageRight := container.NewWithoutLayout()

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
	sortByListOptions := []string{"Country", "Backend", "Backend Score"}
	sortByList := widget.NewSelect(sortByListOptions, func(s string) {
		fmt.Println(s)
	})
	pickRandomButton := widget.NewButton("Random Number", func() {
		numberListLeft.Select(rand.Intn(numberListLeft.Length()))
	})
	horiButtons := container.NewHBox(pickRandomButton, sortByList)

	// TODO: find a better way to stack these two
	leftSide := container.NewVSplit(horiButtons, numberListLeft)
	leftSide.SetOffset(0)
	numberListLeft.Resize(fyne.NewSize(leftSide.Size().Width, w.Canvas().Size().Height-horiButtons.Size().Height))

	// TODO: minimize left side/leading
	split := container.NewHSplit(leftSide, numberPageRight)
	split.SetOffset(0)

	w.SetContent(split)
}

func displayNumberPage(n recvsms.Number, c *fyne.Container) {
	fmt.Println(n)
}

func scrapeNumbers(backends []string, progress *widget.ProgressBar, output *widget.TextGrid) []recvsms.Number {
	defer scrapeInProgress.Unlock()
	scrapeInProgress.Lock()
	progress.Show()
	output.Show()
	nbrChan, beChan := make(chan recvsms.Number, len(backends)), make(chan error)
	for _, backend := range recvsms.ListBackends() {
		go func(backend recvsms.Backend, nbrChan chan recvsms.Number, beChan chan error) {
			nbrs, err := backend.ScrapeNumbers(true)
			if err != nil {
				beChan <- err
				return
			}
			for _, n := range nbrs {
				nbrChan <- n
			}
			beChan <- nil
		}(backend, nbrChan, beChan)
	}
	counter := 0
	numbers := []recvsms.Number{}
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
