package autodiscovery

import(
	"log"
	"os"
	"time"
)

type VolumeState struct {
	dir string
	isConnected bool
	initialDetection bool
	numNewFiles int
	forceFilesStatus bool
}

func GetFileChannel() (channel chan string) {
	channel = make(chan string)
	volumeState := &VolumeState{isConnected: false, initialDetection: false}
	go start(channel, volumeState)
	return
}

func start(channel chan string, volumeState *VolumeState) {
	startDir := detectStartpoint()
	volumeState.setDir(startDir)
	scanner := newFileScaner(startDir)
	for {
		files, err := scanner.GetNewFiles()
		if os.IsNotExist(err) {
			volumeState.setConnected(false)
		} else if err != nil {
			log.Printf("Error while getting new files: %s", err.Error())
		} else {
			volumeState.setConnected(true)
			if volumeState.didInitialDetection() {
				volumeState.setNewFiles(len(*files))
				for _, file := range *files {
					channel <- file
				}
			}
		}
		time.Sleep(3 * time.Second)
	}
}

func (v *VolumeState) setDir(dir string) {
	v.dir = dir
	log.Println("Found new volume: " + dir)
	v.setConnected(true)
}

func (v *VolumeState) setConnected(isConnected bool) {
	if v.isConnected && !isConnected {
		log.Printf("Volume %s is not connected", v.dir)
	} else if !v.isConnected && isConnected {
		v.forceFilesStatus = true
		log.Printf("Volume %s was connected", v.dir)
	}
	v.isConnected = isConnected
}

func (v *VolumeState) didInitialDetection() bool {
	if !v.initialDetection {
		log.Println("Initial files detected. Fly your machines now.")
		v.initialDetection = true
		v.forceFilesStatus = false
		return false
	}
	return true
}

func (v *VolumeState) setNewFiles(num int) {
	if num > 0 || v.forceFilesStatus {
		v.forceFilesStatus = false
		log.Printf("Found %d new files", num)
	}
	v.numNewFiles = num
}
