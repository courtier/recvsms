package main

import (
	"fmt"

	"github.com/courtier/recvsms/pkg/recvsms"
)

func main() {
	for _, backend := range recvsms.ListBackends() {
		nums, err := backend.ScrapeNumbers(false)
		if err != nil {
			panic(err)
		}
		fmt.Println(nums[0])
		msgs, err := backend.ListMessagesForNumber(nums[0], false)
		if err != nil {
			panic(err)
		}
		fmt.Println(msgs[0])
	}
}
