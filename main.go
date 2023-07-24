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
	StartDir        string
	Help            bool
	IncludeAllFiles bool
	HumanReadable   bool
}

func main() {
	fileHashes := make(map[string][]string)

	params, err := getParams()
	if err != nil {
		log.Fatalln(err)
	}

	if params.Help {
		log.Println("you can give one directory path as parameter." +
			"\nwith the option \"-a\" you can include hidden files to the duplicate check." +
			"\nwith the option \"-h\" the output will be formatted human readable.")
		return
	}

	iterateFiles(params.StartDir, fileHashes, params)
	output(fileHashes, params.HumanReadable)
}

func getParams() (Params, error) {
	params := Params{}

	containsParam := func(arg string, param string) bool {
		return strings.HasPrefix(arg, "-") &&
			!strings.HasPrefix(arg, "--") &&
			strings.Contains(arg, param)
	}

	for _, arg := range os.Args[1:] {
		if arg == "--help" {
			params.Help = true
			return params, nil
		}

		if !strings.HasPrefix(arg, "-") && params.StartDir == "" {
			params.StartDir = arg
		} else if !strings.HasPrefix(arg, "-") && params.StartDir != "" {
			err := fmt.Errorf("sorry, but I cannot search in multipe directory paths. can you give me only one directory path, please?")
			return Params{}, err
		}

		if containsParam(arg, "a") {
			params.IncludeAllFiles = true
		}

		if containsParam(arg, "h") {
			params.HumanReadable = true
		}
	}

	if params.StartDir == "" {
		params.StartDir = "."
	}

	return params, nil
}

func iterateFiles(dirPath string, fileHashes map[string][]string, params Params) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatalln(err)
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") && !params.IncludeAllFiles {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())
		if file.IsDir() {
			iterateFiles(filePath, fileHashes, params)
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

func output(fileHashes map[string][]string, humanReadable bool) {
	for hash, filePaths := range fileHashes {
		if len(filePaths) < 2 {
			continue
		}

		if humanReadable {
			log.Printf("=== %s ===\n", hash)
			for _, file := range filePaths {
				log.Println(file)
			}
			log.Print("=========\n\n")
		} else {
			for _, file := range filePaths {
				log.Println(hash, file)
			}
		}
	}
}
