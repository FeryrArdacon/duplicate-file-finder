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
	IncludeAllFiles bool
}

func main() {
	fileHashes := make(map[string][]string)

	if len(os.Args) < 2 {
		log.Fatalln("can you give me one directory path, please? with the option \"-a\" you can include hidden files to the duplicate check.")
	}

	params, err := getParams()
	if err != nil {
		log.Fatalln(err)
	}

	iterateFiles(params.StartDir, fileHashes, params)
	output(fileHashes)
}

func getParams() (Params, error) {
	params := Params{}

	containsParam := func(arg string, param string) bool {
		return strings.HasPrefix(arg, "-") &&
			!strings.HasPrefix(arg, "--") &&
			strings.Contains(arg, param)
	}

	for _, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "-") && params.StartDir == "" {
			params.StartDir = arg
		} else if !strings.HasPrefix(arg, "-") && params.StartDir != "" {
			err := fmt.Errorf("sorry, but I cannot search in multipe directory paths. can you give me only one directory path, please?")
			return Params{}, err
		}

		if containsParam(arg, "a") {
			params.IncludeAllFiles = true
		}
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
			iterateFiles(filepath.Join(filePath), fileHashes, params)
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

func output(fileHashes map[string][]string) {
	for hash, filePaths := range fileHashes {
		if len(filePaths) < 2 {
			continue
		}

		fmt.Printf("=== %s ===\n", hash)
		for _, file := range filePaths {
			fmt.Println(file)
		}
		fmt.Print("=========\n\n")
	}
}
