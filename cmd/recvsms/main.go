package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/courtier/recvsms/pkg/recvsms"
)

var (
	a fyne.App
	w fyne.Window
)

func main() {
	a = app.New()
	w = a.NewWindow("recvsms")
	w.Resize(fyne.NewSize(640, 480))

	backendNames := recvsms.BackendNames()
	backendLength := recvsms.BackendsLength()

	prompt := widget.NewLabel("Pick how you would like the numbers to be scraped:")
	prompt.Alignment = fyne.TextAlignCenter
	prompt.TextStyle.Bold = true

	info := widget.NewLabel(fmt.Sprintf("%d backend(s) available.", backendLength))
	info.Alignment = fyne.TextAlignCenter

	progress := widget.NewProgressBar()
	progress.Min = 0
	progress.Max = float64(len(backendNames))
	progress.Hide()

	output := widget.NewTextGrid()
	output.Hide()

	backendPicker := widget.NewSelect(backendNames, func(opt string) {
		scrapeNumbers([]string{opt}, progress, output)
	})
	backendPicker.PlaceHolder = "Pick a specific backend"

	w.SetContent(container.NewVBox(
		info,
		prompt,
		widget.NewButton("Scrape all backends", func() {
			scrapeNumbers(backendNames, progress, output)
		}),
		backendPicker,
		progress,
		output,
	))

	w.ShowAndRun()
}

func scrapeNumbers(backends []string, progress *widget.ProgressBar, output *widget.TextGrid) []recvsms.Number {
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
				output.SetText(output.Text() + "\nError:" + err.Error())
			}
			if counter == len(backends) {
				break receiveLoop
			}
		}
	}
	return numbers
}

// func main() {
// 	for _, backend := range recvsms.ListBackends() {
// 		nums, err := backend.ScrapeNumbers(false)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(nums[0])
// 		msgs, err := backend.ListMessagesForNumber(nums[0], true)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(msgs[0])
// 		ticker := time.NewTicker(10 * time.Second)
// 		for range ticker.C {
// 			msgs, err = backend.DiffMessagesForNumber(nums[0], true)
// 			if err != nil {
// 				panic(err)
// 			}
// 			fmt.Println(len(msgs))
// 		}
// 	}
// }
