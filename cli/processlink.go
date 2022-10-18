package cli

import (
	"fmt"
	"io/ioutil"
	"log"
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

func DownloadImage(title string, image core.BehanceImage, wg *sync.WaitGroup) error {
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

	// TODO: Create Folder before trying to download
	if _, err = os.Stat(title + "/"); os.IsNotExist(err) {
		_ = os.Mkdir(title, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = ioutil.WriteFile("assets/"+title+"/"+image.FileName, responseBody, 0777)
	if err != nil {
		fmt.Printf("%s \n*******------******* \n", err)
	}

	fmt.Printf("%s Downloaded \n", image.FileName)

	wg.Done()
}
