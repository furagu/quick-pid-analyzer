package autodiscovery

import(
	"time"
	"os"
	"path/filepath"
	"strings"
)

const btlExt = ".bfl"

type FileScaner struct {
	files map[string]int64
	root string
}

func newFileScaner(root string) *FileScaner {
	return &FileScaner{files: make(map[string]int64), root: root}
}

func (s *FileScaner) GetNewFiles() (paths *[]string, err error) {
	if _, err := os.Stat(s.root); err == nil {
		paths = s.checkForFiles()
	} else if os.IsNotExist(err) {
		time.Sleep(2 * time.Second)
	}
	return
}

func (s *FileScaner) checkForFiles() (newFiles *[]string) {
	newFiles = &[]string{}
	err := filepath.Walk(s.root, visit(newFiles, &s.files))
	if err != nil {
		panic(err)
	}
	return
}

func visit(newFiles *[]string, files *map[string]int64) filepath.WalkFunc {
	timestamp := time.Now().Unix()
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if strings.HasPrefix(filepath.Base(path), ".") || strings.ToLower(filepath.Ext(path)) != btlExt || info.IsDir() {
			return nil
		}
		if _, ok := (*files)[path]; !ok {
			(*newFiles) = append(*newFiles, path)
			(*files)[path] = timestamp
		}
		return nil
	}
}
