package autodiscovery

import(
	"time"
	"log"
)


func detectStartpoint() (path string) {
	detected := false
	mountPoints := make(map[string]bool)
	points, err := getMountPoints()
	if err != nil {
		panic(err)
	}
	for _, point := range points {
		mountPoints[point] = true
	}
	log.Print("Please connect you card...")
	for {
		time.Sleep(2 * time.Second)
		log.Print(".")
		points, err := getMountPoints()
		if err != nil {
			panic(err)
		}
		for _, point := range points {
			if _, ok := mountPoints[point]; !ok {
				path = point
				detected = true
				break
			}
		}
		if detected {
			log.Println("")
			break
		}
	}
	return
}
