package cli

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/promisefemi/behance-downloader/core"
)

func ProcessLink(link string) *core.Project {
	projectChan := make(chan *core.Project)
	errorChan := make(chan error)
	// behanceImage := &core.Project{}

	go func() {
		project, err := core.ProcessLink(link)
		if err != nil {
			errorChan <- err
		} else {
			projectChan <- &project
		}
	}()

	for {
		select {
		case behanceImage := <-projectChan:
			fmt.Printf("\r\n")
			return behanceImage
			// break
		case err := <-errorChan:
			fmt.Printf("\r\n")
			exit(fmt.Sprintf("%s", err))
		}
	}

	// fmt.Println("asdfasdh")
	// behanceImage := <-projectChan
	// fmt.Println("asdfasdh")

	// err := <-errorChan

	// if err != nil {
	// 	fmt.Printf("\r\n")
	// 	exit(fmt.Sprintf("%s", err))
	// }

	// if behanceImage.Author.Name != "" {
	// 	fmt.Printf("\r\n")
	// 	wg.Done()
	// 	return behanceImage
	// }

	// return nil

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
	resp, err := http.Get(image.URL)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Cannot reach page, please check URL")
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", destination, image.FileName), responseBody, 0777)
	if err != nil {
		return err
	}

	wg.Done()
	return nil
}
