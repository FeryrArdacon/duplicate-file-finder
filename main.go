package main

import (
	"log"

	"github.com/FeryrArdacon/duplicate-file-finder/filehandler"
	"github.com/FeryrArdacon/duplicate-file-finder/parameters"
)

func main() {
	fileHashes := make(map[string][]string)

	params, err := parameters.GetParams()
	if err != nil {
		log.Fatalln(err)
	}

	if params.Help {
		log.Println("you can give one directory path as parameter." +
			"\nwith the option \"-a\" you can include hidden files to the duplicate check." +
			"\nwith the option \"-h\" the output will be formatted human readable.")
		return
	}

	filehandler.IterateFiles(params.StartDir, fileHashes, params)
	output(fileHashes, params.HumanReadable)
}
