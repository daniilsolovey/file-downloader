package main

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/docopt/docopt-go"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

var (
	version = "0.1"
	usage   = os.ExpandEnv(`
app for downloading addons with .jar extentions to the directory and setting environment variable ADDON_PATH 

<url> 						Downloading from this url
<destination_dir>   		Directory for saving file
<file_name>   				Name of saved file with extention


Usage:
  get-snake-ci-addon   <url>  <destination_dir> <file_name>
  get-snake-ci-addon -h |--help
  get-snake-ci-addon --version

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
		url            = args["<url>"].(string)
		destinationDir = args["<destination_dir>"].(string)
		fileName       = args["<file_name>"].(string)
	)

	pathToFile := filepath.Join(destinationDir, fileName)

	isEmptyDir, err := isEmptyDir(destinationDir)
	if err != nil {
		panic(err)
	}

	if isEmptyDir {
		err = saveAddon(url, destinationDir, fileName)
		if err != nil {
			panic(err)
		}

		absPathToCreatedFile, err := filepath.Abs(pathToFile)
		if err != nil {
			panic(err)
		}

		err = os.Setenv("ADDON_PATH", absPathToCreatedFile)
		if err != nil {
			panic(err)
		}

		log.Infof(
			nil,
			"addon saved, environment variable: ADDON_PATH:%s",
			absPathToCreatedFile,
		)

		return
	}

	fileInDir, err := getNameOfOneFileInDirectory(destinationDir)
	if err != nil {
		panic(err)
	}

	pathToExistsFile := filepath.Join(destinationDir, fileInDir)
	absPathToExistsFile, err := filepath.Abs(pathToExistsFile)
	if err != nil {
		panic(err)
	}

	os.Setenv("ADDON_PATH", absPathToExistsFile)
	if err != nil {
		panic(err)
	}

	log.Infof(
		nil,
		"addon already exists, environment variable: ADDON_PATH:%s",
		absPathToExistsFile,
	)

}

func saveAddon(url, destinationDir, fileName string) error {
	log.Infof(nil, "get addon by url: %s", url)
	response, err := http.Get(url)
	if err != nil {
		return karma.Format(
			err,
			"unable to get response by url: %s",
			url,
		)
	}

	defer response.Body.Close()

	isNotExistDir := false
	if _, err := os.Stat(destinationDir); os.IsNotExist(err) {
		isNotExistDir = true
	}

	if isNotExistDir {
		err := os.MkdirAll(destinationDir, os.ModePerm)
		if err != nil {
			return karma.Format(
				err,
				"unable to create directory by path: %s",
				destinationDir,
			)
		}

	}

	filePath := filepath.Join(destinationDir, fileName)
	out, err := os.Create(filePath)
	if err != nil {
		if err != nil {
			return karma.Format(
				err,
				"unable to create file with name: %s",
				filePath,
			)
		}
	}

	defer out.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		if err != nil {
			return karma.Format(
				err,
				"unable to write response from url to file, url:%s, file: %s",
				url,
				filePath,
			)
		}
	}

	return nil
}

func isEmptyDir(directory string) (bool, error) {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return true, nil
	}

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return false, err
	}

	if len(files) == 0 {
		return true, nil
	}

	return false, nil
}

func getNameOfOneFileInDirectory(directory string) (string, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		err := errors.New("directory is empty")
		return "", err
	}

	if len(files) > 1 {
		err := errors.New("this direcory contains more than one file")
		return "", err
	}

	path := filepath.Join(directory, files[0].Name())
	if filepath.Ext(path) != ".jar" {
		err := errors.New("one file exists in directory, but this file doesn't have .jar extention")
		return "", err
	}

	return files[0].Name(), nil
}
