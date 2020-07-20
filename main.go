package main

import (
	"io"
	"net/http"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/reconquest/pkg/log"
)

var (
	version = "0.1"
	usage   = os.ExpandEnv(`
app for downloading snake-ci plugin .jar file 

<url> 						Downloading from this url
<file_name>   				Name of saved file with extention

Usage:
  get-snake-ci-addon   <url>  <file_name>

  stacket -h |--help
  stacket --version

Options:
  -h --help                     Show this screen.
  --version                     Show version.
`)
)

func main() {
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		panic(err)
	}

	var (
		url      = args["<url>"].(string)
		fileName = args["<file_name>"].(string)
	)

	log.Infof(nil, "get addon by url: %s", url)
	err = saveAddon(url, fileName)
	if err != nil {
		panic(err)
	}
}

func saveAddon(url, filepath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return err
	}

	return nil
}
