package main

import (
	"log"
	"github.com/furagu/quick-pid-analyzer/src/autodiscovery"
)

func main() {
	newFiles := autodiscovery.GetFileChannel()

	for filePath := range newFiles {
		log.Println("New file: " + filePath)
	}
}
