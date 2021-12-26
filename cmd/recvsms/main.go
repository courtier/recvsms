package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/courtier/recvsms/pkg/recvsms"
)

func main() {
	a := app.New()
	w := a.NewWindow("recvsms")
	s := fyne.NewSize(640, 480)
	w.Resize(s)
	prompt := widget.NewLabel("Pick how you would like the numbers to be scraped.")
	prompt.Alignment = fyne.TextAlignCenter
	info := widget.NewLabel(fmt.Sprintf("%d backend(s) available.", recvsms.BackendsLength()))
	info.Alignment = fyne.TextAlignCenter
	w.SetContent(container.NewVBox(
		prompt,
		info,
		widget.NewButton("Scrape All Numbers From All Backends", func() {
			prompt.SetText("Welcome :)")
		}),
		widget.NewButton("Hi!", func() {
			prompt.SetText("Welcome :)")
		}),
	))
	w.ShowAndRun()
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
