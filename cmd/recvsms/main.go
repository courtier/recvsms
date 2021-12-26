package main

import (
	"fmt"
	"time"

	"github.com/courtier/recvsms/pkg/recvsms"
)

func main() {
	for _, backend := range recvsms.ListBackends() {
		nums, err := backend.ScrapeNumbers(false)
		if err != nil {
			panic(err)
		}
		fmt.Println(nums[0])
		msgs, err := backend.ListMessagesForNumber(nums[0], true)
		if err != nil {
			panic(err)
		}
		fmt.Println(msgs[0])
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			msgs, err = backend.DiffMessagesForNumber(nums[0], true)
			if err != nil {
				panic(err)
			}
			fmt.Println(len(msgs))
		}
	}
}
