package main

import (
	"github.com/courtier/recvsms/pkg/recvsms"
)

func main() {
	print("hello")
	for _, backend := range recvsms.ListBackends() {
		print(backend)
	}
}
