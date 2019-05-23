package autodiscovery

import(
	"time"
	"os"
	"log"
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

func (s *FileScaner) GetNewFiles() (newFiles *[]string, err error) {
	log.Println("Checking for new files")
	newFiles = &[]string{}
	err = filepath.Walk(s.root, visit(newFiles, &s.files))
	if err != nil {
		return
	}
	log.Printf("Found %d new files", len(*newFiles))
	return
}

func visit(newFiles *[]string, files *map[string]int64) filepath.WalkFunc {
	timestamp := time.Now().Unix()
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
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
