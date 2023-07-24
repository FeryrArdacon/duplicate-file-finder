package filehandler

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/FeryrArdacon/duplicate-file-finder/parameters"
)

func IterateFiles(dirPath string, fileHashes map[string][]string, params parameters.Params) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatalln(err)
	}

	for _, file := range files {
		isHidden, err := isHiddenFile(file.Name(), dirPath)
		if err != nil {
			log.Fatalln(err)
		}

		if isHidden && !params.IncludeAllFiles {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())
		if file.IsDir() {
			IterateFiles(filePath, fileHashes, params)
			continue
		}

		hash := fmt.Sprintf("%x", getHashForFile(filePath))
		fileHashes[hash] = append(fileHashes[hash], filePath)
	}
}

func getHashForFile(filePath string) []byte {
	fileForHashing, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer fileForHashing.Close()

	hashSha256 := sha256.New()
	if _, err := io.Copy(hashSha256, fileForHashing); err != nil {
		log.Fatalln(err)
	}

	return hashSha256.Sum(nil)
}
