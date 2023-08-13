package main

import (
	"fmt"
	"os"
)

const (
	maxFile = 1000
)

func main() {
	for i := 0; i < maxFile; i++ {
		file, err := os.Create(fmt.Sprintf("./temp/%v-name.text", i))
		if err != nil {
			panic(err)
		}
		defer file.Close()

		file.WriteString(fmt.Sprintf("%v of %v", i, maxFile))
	}
}
