package autodiscovery

import(
	"log"
	"os"
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
	initialFilesDetected := false
	notConnectedReported := false
	for {
		files, err := scanner.GetNewFiles()
		if os.IsNotExist(err) {
			if !notConnectedReported {
				notConnectedReported = true
				log.Printf("Volume %s is not connected", startDir)
			}
		} else if err != nil {
			log.Printf("Error while getting new files: %s", err.Error())
		} else {
			notConnectedReported = false
			if !initialFilesDetected {
				initialFilesDetected = true
				log.Println("Initial files detected. Fly your machines now.")
			} else {
				log.Printf("Found %d new files", len(*files))
				for _, file := range *files {
					channel <- file
				}
			}
		}
		time.Sleep(3 * time.Second)
	}
}
