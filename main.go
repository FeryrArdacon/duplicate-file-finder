package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Params struct {
	IncludeAllFiles bool
}

func main() {
	fileHashes := make(map[string][]string)

	if len(os.Args) < 2 {
		log.Fatalln("Can you give me one directory path, please? With the option \"-a\" you can include hidden files to the duplicate check.")
		return
	}

	var params = Params{}
	var startDir = ""

	for _, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "-") && startDir == "" {
			startDir = arg
		} else if !strings.HasPrefix(arg, "-") && startDir != "" {
			log.Fatalln("Sorry, but I cannot search in multipe directory paths. Can you give me only one directory path, please?")
			return
		}

		if strings.HasPrefix(arg, "-") && strings.Contains(arg, "a") {
			params.IncludeAllFiles = true
		}
	}

	iterateFiles(startDir, fileHashes, params)

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

func iterateFiles(dirPath string, fileHashes map[string][]string, params Params) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") && !params.IncludeAllFiles {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())
		if file.IsDir() {
			iterateFiles(filepath.Join(filePath), fileHashes, params)
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
