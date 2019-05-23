package autodiscovery

import(
	"log"
	"time"
)

func GetFileChannel() (channel chan string) {
	channel = make(chan string)
	go start(channel)
	return
}

func start(channel chan string) {
	startDir := detectStartpoint()
	log.Println("Found new volume: " + startDir)
	scanner := newFileScaner(startDir)
	for {
		files, err := scanner.GetNewFiles()
		if err != nil {
			panic(err)
		}
		for _, file := range *files {
			channel <- file
		}
		time.Sleep(3 * time.Second)
	}
}
