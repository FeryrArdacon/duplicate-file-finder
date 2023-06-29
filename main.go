package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	fileHashes := make(map[string][]string)

	iterateFiles(".", fileHashes)

	for hash, filePaths := range fileHashes {
		if len(filePaths) > 1 {
			fmt.Printf("=== %s ===\n", hash)
			for _, file := range filePaths {
				fmt.Println(file)
			}
			fmt.Print("=========\n\n")
		}
	}
}

func iterateFiles(dirPath string, fileHashes map[string][]string) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		filePath := filepath.Join(dirPath, file.Name())
		if file.IsDir() {
			iterateFiles(filepath.Join(filePath), fileHashes)
		} else {
			hash := fmt.Sprintf("%x", getHashForFile(filePath))
			fileHashes[hash] = append(fileHashes[hash], filePath)
		}
	}
}

func getHashForFile(filePath string) []byte {
	fileForHashing, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer fileForHashing.Close()

	hashSha256 := sha256.New()
	if _, err := io.Copy(hashSha256, fileForHashing); err != nil {
		log.Fatal(err)
	}

	return hashSha256.Sum(nil)
}
