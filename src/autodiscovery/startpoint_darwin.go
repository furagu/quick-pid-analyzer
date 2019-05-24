// +build darwin

package autodiscovery

import(
	"os"
)

const volumesDir = "/Volumes"

func getMountPoints() (points []string, err error) {
	f, err := os.Open(volumesDir)
	defer f.Close()
    if err != nil {
		return
    }
	list, err := f.Readdir(-1)
	if err != nil {
		return
	}
	for _, item := range list {
		if item.IsDir() {
			points = append(points, volumesDir + "/" + item.Name())
		}
	}
	return
}
