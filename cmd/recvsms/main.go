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
		fmt.Println(nums)
		fmt.Println(len(nums))
	}
}
