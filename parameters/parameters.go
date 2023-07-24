package parameters

import (
	"fmt"
	"os"
	"strings"
)

type Params struct {
	StartDir        string
	Help            bool
	IncludeAllFiles bool
	HumanReadable   bool
}

func GetParams() (Params, error) {
	params := Params{}

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

func containsParam(arg string, param string) bool {
	return strings.HasPrefix(arg, "-") &&
		!strings.HasPrefix(arg, "--") &&
		strings.Contains(arg, param)
}
