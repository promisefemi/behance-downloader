package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/promisefemi/behance-downloader/core"
)

func ProcessLink(link string) (*core.Project, error) {
	project, err := core.ProcessLink(link)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func ListDetails(behanceImage *core.Project) {
	PrintAuthor(behanceImage.Author, false)
	for i, image := range behanceImage.Images {
		fmt.Printf("%d - %s\n", i+1, image.URL)
	}
}

func PrintAuthor(author core.Author, onlyAuthor bool) {

	fmt.Printf("%s -- %s (%s)\n", "Author", author.Name, author.ProfileLink)
	if onlyAuthor {
		fmt.Printf("\n")
	}
}

func CreateFolder(folder string) error {
	if _, err := os.Stat(folder + "/"); os.IsNotExist(err) {
		err = os.MkdirAll(folder, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func DownloadImage(destination string, image core.BehanceImage, wg *sync.WaitGroup) error {
	// fmt.Printf("URL of the Image %s \n", imageURL)
	body, _, _, err := core.ProcessDownload(image.URL)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", destination, image.FileName), body, 0777)
	if err != nil {
		return err
	}

	wg.Done()
	return nil
}
