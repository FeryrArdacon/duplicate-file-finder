package main

import "log"

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
